// Copyright 2023 Talhuang<talhuang1231@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package server

import (
	"net"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	homedir "github.com/skeleton1231/gotal/pkg/util/common"
	"github.com/spf13/viper"
)

// ProjectName defines the name of this API server project.
// Constants for default configurations.
const (
	ProjectName          = "APISERVER"
	DefaultConfigType    = "yaml"
	DefaultEnvKeyReplace = "_"
)

// DefaultConfig contains default values for server configurations.
type DefaultConfig struct {
	HomeDir    string
	EnvPrefix  string
	ConfigPath string
	JwtRealm   string
}

var defaultConf = &DefaultConfig{
	HomeDir:    "." + strings.ToLower(ProjectName),
	EnvPrefix:  strings.ToUpper(ProjectName),
	ConfigPath: "/etc/" + strings.ToLower(ProjectName),
	JwtRealm:   strings.ToLower(ProjectName) + " jwt",
}

// Config represents the main configuration structure for the API server.
type Config struct {
	SecureServing   *SecureServingInfo
	InsecureServing *InsecureServingInfo
	Jwt             *JwtInfo
	Mode            string
	Middlewares     []string
	Healthz         bool
	EnableProfiling bool
	EnableMetrics   bool
	RateLimit       *RateLimit
}

// CertKey represents the certificate and key configuration for secure serving.
type CertKey struct {
	CertFile string // Path to the PEM-encoded certificate.
	KeyFile  string // Path to the PEM-encoded private key associated with the certificate.
}

// SecureServingInfo encapsulates the configuration for serving the API over HTTPS.
type SecureServingInfo struct {
	BindAddress string
	BindPort    int
	CertKey     CertKey
}

// Address constructs a complete address by combining BindAddress and BindPort.
func (s *SecureServingInfo) Address() string {
	return net.JoinHostPort(s.BindAddress, strconv.Itoa(s.BindPort))
}

// InsecureServingInfo contains the configuration for serving the API over HTTP.
type InsecureServingInfo struct {
	Address string
}

// JwtInfo defines configuration parameters for JWT-based authentication.
type JwtInfo struct {
	Realm      string
	Key        string
	Timeout    time.Duration
	MaxRefresh time.Duration
}

// RateLimitConfig represents the configuration for rate limiting.
type RateLimit struct {
	TokensPerSecond int  // Number of tokens generated per second.
	Burst           int  // Maximum burst size.
	Enabled         bool // Enable rate limiting.
}

// NewConfig creates and returns a new Config instance with default settings.
func NewConfig() *Config {
	return &Config{
		Healthz:         true,
		Mode:            gin.ReleaseMode,
		Middlewares:     []string{},
		EnableProfiling: true,
		EnableMetrics:   true,
		Jwt: &JwtInfo{
			Realm:      defaultConf.JwtRealm,
			Timeout:    1 * time.Hour,
			MaxRefresh: 1 * time.Hour,
		},
		// Default rate limit configurations
		RateLimit: &RateLimit{
			TokensPerSecond: 1,    // Example default value
			Burst:           10,   // Example default value
			Enabled:         true, // Enable by default
		},
	}
}

// CompletedConfig represents a configuration that has been finalized and is ready for use.
type CompletedConfig struct {
	*Config
}

// Complete finalizes the Config by setting any missing values and returns a CompletedConfig.
func (c *Config) Complete() CompletedConfig {
	return CompletedConfig{c}
}

// New initializes and returns a new APIServer instance based on the completed configuration.
func (c CompletedConfig) New() (*APIServer, error) {
	// Set the operational mode for the Gin web framework.
	gin.SetMode(c.Mode)

	// Create a new API server instance with the specified configurations.
	s := &APIServer{
		SecureServingInfo:   c.SecureServing,
		InsecureServingInfo: c.InsecureServing,
		healthz:             c.Healthz,
		enableMetrics:       c.EnableMetrics,
		enableProfiling:     c.EnableProfiling,
		middlewares:         c.Middlewares,
		Engine:              gin.New(),
		ShutdownTimeout:     30 * time.Second,
	}

	// Initialize the API server with the required setup.
	initAPIServer(s)

	return s, nil
}

// LoadConfig loads the API server configuration from a given file or a default location.
func LoadConfig(cfg string, defaultName string) error {
	if cfg != "" {
		viper.SetConfigFile(cfg)
	} else {
		viper.AddConfigPath(".")
		viper.AddConfigPath(filepath.Join(homedir.HomeDir(), defaultConf.HomeDir))
		viper.AddConfigPath(defaultConf.ConfigPath)
		viper.SetConfigName(defaultName)
	}

	viper.SetConfigType(DefaultConfigType)
	viper.AutomaticEnv()
	viper.SetEnvPrefix(defaultConf.EnvPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", DefaultEnvKeyReplace, "-", DefaultEnvKeyReplace))

	if err := viper.ReadInConfig(); err != nil {
		logrus.Warnf("WARNING: viper failed to discover and load the configuration file: %s", err.Error())
		return err
	}
	return nil
}
