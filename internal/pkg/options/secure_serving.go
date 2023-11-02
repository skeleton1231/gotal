// Copyright 2023 Talhuang<talhuang1231@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package options

import (
	"fmt"
	"net"
	"path"

	"github.com/skeleton1231/gotal/internal/pkg/server"
	"github.com/spf13/pflag"
)

// Default values for secure serving options
const (
	DefaultBindAddress   = "0.0.0.0"
	DefaultBindPort      = 8443
	DefaultPairName      = "apiserver"
	DefaultCertDirectory = "/var/run/apiserver"
)

// SecureServingOptions represents configuration for securely serving requests (usually over HTTPS).
type SecureServingOptions struct {
	BindAddress string `json:"bind-address" mapstructure:"bind-address"`
	BindPort    int    `json:"bind-port"    mapstructure:"bind-port"`
	Required    bool
	ServerCert  GeneratableKeyCert `json:"tls"          mapstructure:"tls"`
}

// CertKey holds paths for a certificate and its corresponding key file.
type CertKey struct {
	CertFile string `json:"cert-file"        mapstructure:"cert-file"`
	KeyFile  string `json:"private-key-file" mapstructure:"private-key-file"`
}

// GeneratableKeyCert holds configuration for generating or locating key and cert files.
type GeneratableKeyCert struct {
	CertKey CertKey `json:"cert-key" mapstructure:"cert-key"`

	CertDirectory string `json:"cert-dir"  mapstructure:"cert-dir"`

	PairName string `json:"pair-name" mapstructure:"pair-name"`
}

// NewSecureServingOptions creates a new instance of SecureServingOptions with default values.
func NewSecureServingOptions() *SecureServingOptions {
	return &SecureServingOptions{
		BindAddress: DefaultBindAddress,
		BindPort:    DefaultBindPort,
		Required:    true,
		ServerCert: GeneratableKeyCert{
			PairName:      DefaultPairName,
			CertDirectory: DefaultCertDirectory,
		},
	}
}

// ApplyTo sets up the secure serving info on the given server Config.
func (s *SecureServingOptions) ApplyTo(c *server.Config) error {
	c.SecureServing = &server.SecureServingInfo{
		BindAddress: s.BindAddress,
		BindPort:    s.BindPort,
		CertKey: server.CertKey{
			CertFile: s.ServerCert.CertKey.CertFile,
			KeyFile:  s.ServerCert.CertKey.KeyFile,
		},
	}

	return nil
}

// Validate checks for any errors in the secure serving options.
func (s *SecureServingOptions) Validate() []error {
	if s == nil {
		return nil
	}

	errors := []error{}

	if s.Required && s.BindPort < 1 || s.BindPort > 65535 {
		errors = append(
			errors,
			fmt.Errorf(
				"--secure.bind-port %v must be between 1 and 65535, inclusive. It cannot be turned off with 0",
				s.BindPort,
			),
		)
	} else if s.BindPort < 0 || s.BindPort > 65535 {
		errors = append(errors, fmt.Errorf("--secure.bind-port %v must be between 0 and 65535, inclusive. 0 for turning off secure port", s.BindPort))
	}

	return errors
}

// AddFlags adds flags for the secure serving options to the given flag set.
func (s *SecureServingOptions) AddFlags(fs *pflag.FlagSet) {
	// IP address configuration for secure port
	fs.StringVar(&s.BindAddress, "secure.bind-address", s.BindAddress,
		"IP address to listen on for the secure port. Set to '0.0.0.0' for all IPv4 or '::' for all IPv6 interfaces. If left empty, all available interfaces are used.")

	// Description for secure port configuration
	desc := "Port for serving HTTPS with necessary authentication and authorization."
	if s.Required {
		desc += " This port cannot be disabled by setting it to 0."
	} else {
		desc += " Set to 0 to disable HTTPS."
	}
	fs.IntVar(&s.BindPort, "secure.bind-port", s.BindPort, desc)

	// Directory containing TLS certificates
	fs.StringVar(&s.ServerCert.CertDirectory, "secure.tls.cert-dir", s.ServerCert.CertDirectory,
		"Directory holding the TLS certificates. If specific cert and private key files are provided, this directory setting is ignored.")

	// Naming convention for the TLS pair files
	fs.StringVar(&s.ServerCert.PairName, "secure.tls.pair-name", s.ServerCert.PairName,
		"Naming convention used alongside the cert directory to determine the cert and key file names. Results in filenames like '<cert-dir>/<pair-name>.crt' and '<cert-dir>/<pair-name>.key'.")

	// File containing the default certificate for HTTPS
	fs.StringVar(&s.ServerCert.CertKey.CertFile, "secure.tls.cert-key.cert-file", s.ServerCert.CertKey.CertFile,
		"File holding the default x509 certificate for HTTPS. If a CA certificate exists, it should be concatenated after the server certificate in this file.")

	// File containing the default private key
	fs.StringVar(&s.ServerCert.CertKey.KeyFile, "secure.tls.cert-key.private-key-file", s.ServerCert.CertKey.KeyFile,
		"File holding the x509 private key that corresponds to the certificate in --secure.tls.cert-key.cert-file.")
}

// Complete finalizes the SecureServingOptions after flag parsing.
func (s *SecureServingOptions) Complete() error {
	if s == nil || s.BindPort == 0 {
		return nil
	}

	keyCert := &s.ServerCert.CertKey
	if len(keyCert.CertFile) != 0 || len(keyCert.KeyFile) != 0 {
		return nil
	}

	if len(s.ServerCert.CertDirectory) > 0 {
		if len(s.ServerCert.PairName) == 0 {
			return fmt.Errorf("--secure.tls.pair-name is required if --secure.tls.cert-dir is set")
		}
		keyCert.CertFile = path.Join(s.ServerCert.CertDirectory, s.ServerCert.PairName+".crt")
		keyCert.KeyFile = path.Join(s.ServerCert.CertDirectory, s.ServerCert.PairName+".key")
	}

	return nil
}

// CreateListener sets up a TCP listener on the specified address.
func CreateListener(addr string) (net.Listener, int, error) {
	network := "tcp"

	ln, err := net.Listen(network, addr)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to listen on %v: %w", addr, err)
	}

	// get port
	tcpAddr, ok := ln.Addr().(*net.TCPAddr)
	if !ok {
		_ = ln.Close()

		return nil, 0, fmt.Errorf("invalid listen address: %q", ln.Addr().String())
	}

	return ln, tcpAddr.Port, nil
}
