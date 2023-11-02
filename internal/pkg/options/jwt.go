// Copyright 2023 Talhuang<talhuang1231@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package options

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/skeleton1231/gotal/internal/pkg/server"
	"github.com/spf13/pflag"
)

// JwtOptions defines the configuration options related to JWT tokens.
type JwtOptions struct {
	Realm      string        `json:"realm"       mapstructure:"realm" validate:"required"`
	Key        string        `json:"key"         mapstructure:"key" validate:"required,len=32"`
	Timeout    time.Duration `json:"timeout"     mapstructure:"timeout" validate:"required"`
	MaxRefresh time.Duration `json:"max-refresh" mapstructure:"max-refresh" validate:"required"`
}

// NewJwtOptions returns a JwtOptions object with default values.
func NewJwtOptions() *JwtOptions {
	defaults := server.NewConfig()
	return &JwtOptions{
		Realm:      defaults.Jwt.Realm,
		Key:        defaults.Jwt.Key,
		Timeout:    defaults.Jwt.Timeout,
		MaxRefresh: defaults.Jwt.MaxRefresh,
	}
}

// ApplyTo applies the current options to the provided server configuration.
func (s *JwtOptions) ApplyTo(c *server.Config) error {
	c.Jwt = &server.JwtInfo{
		Realm:      s.Realm,
		Key:        s.Key,
		Timeout:    s.Timeout,
		MaxRefresh: s.MaxRefresh,
	}
	return nil
}

// Validate checks and validates the user-provided parameters during program startup.
func (s *JwtOptions) Validate() []error {
	var errs []error

	err := validate.Struct(s)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			errs = append(errs, err)
		}

		for _, err := range err.(validator.ValidationErrors) {
			errs = append(errs, fmt.Errorf("%s in field %s must be valid", err.Tag(), err.Field()))
		}
	}
	return errs
}

// AddFlags adds the JWT-related flags to the provided FlagSet.
func (s *JwtOptions) AddFlags(fs *pflag.FlagSet) {
	if fs == nil {
		return
	}

	// Adding command-line flags related to JWT configuration.
	fs.StringVar(&s.Realm, "jwt.realm", s.Realm, "Set the realm name displayed to the user.")
	fs.StringVar(&s.Key, "jwt.key", s.Key, "Set the private key used for signing the JWT token.")
	fs.DurationVar(&s.Timeout, "jwt.timeout", s.Timeout, "Set the expiration duration for the JWT token.")
	fs.DurationVar(&s.MaxRefresh, "jwt.max-refresh", s.MaxRefresh, "Set the maximum duration for token refresh by clients.")
}
