// Copyright 2023 Talhuang<talhuang1231@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package options provides configurable options for server setup and runtime.
package options

import (
	"fmt"
	"net"
	"strconv"

	"github.com/skeleton1231/gotal/internal/pkg/server"
	"github.com/spf13/pflag"
)

// InsecureServingOptions represents configuration options for insecure server communication (HTTP).
type InsecureServingOptions struct {
	BindAddress string `json:"bind-address" mapstructure:"bind-address"` // IP address to bind the server.
	BindPort    int    `json:"bind-port"    mapstructure:"bind-port"`    // Port number to bind the server.
}

// NewInsecureServingOptions initializes a new InsecureServingOptions object with default values.
func NewInsecureServingOptions() *InsecureServingOptions {
	return &InsecureServingOptions{
		BindAddress: "127.0.0.1", // Default bind address is set to localhost.
		BindPort:    8080,        // Default port is set to 8080.
	}
}

// ApplyTo updates the provided server.Config with the settings from InsecureServingOptions.
func (s *InsecureServingOptions) ApplyTo(c *server.Config) error {
	c.InsecureServing = &server.InsecureServingInfo{
		Address: net.JoinHostPort(s.BindAddress, strconv.Itoa(s.BindPort)), // Combine IP and Port to a single address string.
	}

	return nil
}

// Validate checks if the settings of InsecureServingOptions are valid.
func (s *InsecureServingOptions) Validate() []error {
	var errors []error

	// Ensure the bind port is within the valid range.
	if s.BindPort < 0 || s.BindPort > 65535 {
		errors = append(
			errors,
			fmt.Errorf(
				"--insecure.bind-port %v must be between 0 and 65535, inclusive. 0 for turning off insecure (HTTP) port",
				s.BindPort,
			),
		)
	}

	return errors
}

// AddFlags adds the InsecureServingOptions flags to the provided FlagSet.
func (s *InsecureServingOptions) AddFlags(fs *pflag.FlagSet) {
	// Define the bind address flag.
	fs.StringVar(&s.BindAddress, "insecure.bind-address", s.BindAddress,
		"Specifies the IP address on which the service will listen. Use '0.0.0.0' for all IPv4 interfaces and '::' for all IPv6 interfaces.")
	// Define the bind port flag.
	fs.IntVar(&s.BindPort, "insecure.bind-port", s.BindPort,
		"Specifies the port for serving unsecured and unauthenticated access. Ensure firewall settings prevent external access to this port. In default configurations, port 443 on the public address is proxied to this port by nginx. Set to zero to disable this port.")
}
