// Copyright 2023 Talhuang<talhuang1231@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	homedir "github.com/skeleton1231/gotal/pkg/util/common"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Name of the flag for configuration file.
const configFlagName = "config"

// Variable to store the path of the configuration file provided by the user.
var cfgFile string

// init function is executed when the package is imported.
func init() {
	// Define a command-line flag named "config" (short form "c")
	// to specify the configuration file.
	pflag.StringVarP(&cfgFile, "config", "c", cfgFile, "Specify a configuration file (`FILE`). Supports JSON, TOML, YAML properties formats.")

}

// addConfigFlag configures viper to use a config file and sets up the environment variables.
func addConfigFlag(basename string, fs *pflag.FlagSet) {
	// Add the configuration flag to the provided flag set.
	fs.AddFlag(pflag.Lookup(configFlagName))

	// Configure viper to read configuration values from environment variables.
	viper.AutomaticEnv()
	// Set the environment variable prefix based on the application basename.
	viper.SetEnvPrefix(strings.Replace(strings.ToUpper(basename), "-", "_", -1))
	// Replace certain characters in environment variable names.
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// Configure viper with the path and name of the configuration file.
	cobra.OnInitialize(func() {
		// If a configuration file path was provided by the user, use it.
		if cfgFile != "" {
			viper.SetConfigFile(cfgFile)
		} else {
			// Otherwise, set default paths where viper should look for the configuration file.
			viper.AddConfigPath(".")
			// If the basename contains a hyphen, determine potential config directories based on the first part.
			if names := strings.Split(basename, "-"); len(names) > 1 {
				viper.AddConfigPath(filepath.Join(homedir.HomeDir(), "."+names[0]))
				viper.AddConfigPath(filepath.Join("/etc", names[0]))
			}
			// Set the name of the configuration file (without extension).
			viper.SetConfigName(basename)
		}

		// Try to read the configuration file.
		if err := viper.ReadInConfig(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error: failed to read configuration file(%s): %v\n", cfgFile, err)
			os.Exit(1)
		}
	})
}

// printConfig prints the loaded configuration values to the console.
func printConfig() {
	// Fetch all configuration keys from viper.
	if keys := viper.AllKeys(); len(keys) > 0 {
		// This appears to have a missing definition for 'progressMessage'.
		// Assuming it's a global string that indicates some progress state.
		fmt.Printf("%v Configuration items:\n", progressMessage)
		// Loop through and print each key-value pair.
		for _, k := range keys {
			fmt.Printf("%s: %v\n", k, viper.Get(k))
		}
	}
}
