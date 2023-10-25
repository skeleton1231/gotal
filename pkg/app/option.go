// Copyright 2023 Talhuang<talhuang1231@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package app

import (
	sections "github.com/skeleton1231/gotal/pkg/util/flag"
	"github.com/spf13/viper"
)

// CliOptions defines an interface for command-line options.
// This interface provides methods for getting flag sets and validating those flags.
type CliOptions interface {
	// Flags returns a set of named flag sets that can be used by the application.
	// NamedFlagSets helps organize flags into logical groups for better CLI display.
	Flags() (fss sections.NamedFlagSets)

	// Validate checks the validity of the provided command-line options.
	// It returns a list of errors encountered during validation.
	Validate() []error
}

// ConfigurableOptions defines an interface for options that can be configured using viper.
// Viper is a configuration solution for Go applications that supports setting defaults, reading
// from JSON, TOML, YAML, and more.
type ConfigurableOptions interface {
	// ApplyFlags applies the values from the provided viper instance to the implementing struct.
	// This allows the command-line flags to override values from other sources such as configuration files.
	ApplyFlags(v *viper.Viper) []error
}

// CompleteableOptions defines an interface for options that require completion.
// Some options might need additional setup or data fetching after being parsed and validated.
type CompleteableOptions interface {
	// Complete completes any additional setup or data fetching required for the option.
	Complete() error
}

// PrintableOptions defines an interface for options that can be converted to a string.
// This interface provides a method to get a string representation of the option,
// which can be useful for logging or display purposes.
type PrintableOptions interface {
	// String returns a string representation of the option.
	String() string
}
