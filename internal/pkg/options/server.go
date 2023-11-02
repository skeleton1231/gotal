// Copyright 2023 Talhuang<talhuang1231@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package options

import (
	"github.com/go-playground/validator/v10"
	"github.com/skeleton1231/gotal/internal/pkg/server"
	"github.com/spf13/pflag"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ServerRunOptions represents the customizable options for running the server.
// It includes fields such as Mode, Healthz status, and the middlewares used.
type ServerRunOptions struct {
	Mode        string   `json:"mode"        mapstructure:"mode"`        // Mode specifies the server's operating mode (e.g., debug, test, release).
	Healthz     bool     `json:"healthz"     mapstructure:"healthz"`     // Healthz indicates whether a /healthz endpoint should be established for readiness checks.
	Middlewares []string `json:"middlewares" mapstructure:"middlewares"` // Middlewares lists the middlewares allowed for the server.
}

// NewServerRunOptions initializes a new ServerRunOptions instance with default settings from the server's configuration.
func NewServerRunOptions() *ServerRunOptions {
	defaults := server.NewConfig() // Fetch default configurations.

	return &ServerRunOptions{ // Populate the ServerRunOptions with default values.
		Mode:        defaults.Mode,
		Healthz:     defaults.Healthz,
		Middlewares: defaults.Middlewares,
	}
}

// ApplyTo updates the given server configuration (c) with the values from the ServerRunOptions.
func (s *ServerRunOptions) ApplyTo(c *server.Config) error {
	c.Mode = s.Mode
	c.Healthz = s.Healthz
	c.Middlewares = s.Middlewares

	return nil // Return nil as there's no error handling currently.
}

// Validate checks the ServerRunOptions for any inconsistencies or errors.
// Currently, it always returns an empty error slice but can be extended for more complex validations.
func (s *ServerRunOptions) Validate() []error {
	errors := []error{}

	return errors
}

// AddFlags binds the ServerRunOptions fields to command-line flags using the given FlagSet (fs).
// This allows users to customize server behavior via command-line options.
func (s *ServerRunOptions) AddFlags(fs *pflag.FlagSet) {
	// Bind Mode field to --server.mode flag.
	fs.StringVar(&s.Mode, "server.mode", s.Mode,
		"Define the operating mode of the server. Supported modes include: debug, test, and release.")

	// Bind Healthz field to --server.healthz flag.
	fs.BoolVar(&s.Healthz, "server.healthz", s.Healthz,
		"Enable a self-readiness check and establish a /healthz endpoint.")

	// Bind Middlewares field to --server.middlewares flag.
	fs.StringSliceVar(&s.Middlewares, "server.middlewares", s.Middlewares,
		"Specify a comma-separated list of allowed middlewares for the server. Defaults will be used if this list is empty.")
}
