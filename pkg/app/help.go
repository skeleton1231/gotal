// Copyright 2023 Talhuang<talhuang1231@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package app

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	// flagHelp defines the long version of the help flag.
	flagHelp = "help"

	// flagHelpShorthand defines the short version of the help flag.
	flagHelpShorthand = "H"
)

// helpCommand returns a new cobra.Command configured to show help for other commands.
// It allows the user to get detailed information about individual commands.
func helpCommand(name string) *cobra.Command {
	return &cobra.Command{
		Use:   "help [command]",
		Short: "Display help information about a command.",
		Long: fmt.Sprintf(`Use this command to get detailed help information about any other command.
For example, to get help for a specific command, type: %s help [command name]`, name),
		Run: func(c *cobra.Command, args []string) {
			// Find the command based on the provided arguments.
			cmd, _, e := c.Root().Find(args)
			if cmd == nil || e != nil {
				// If the command doesn't exist, print an unknown topic message.
				c.Printf("Unknown help topic: %#q\n", args)
				_ = c.Root().Usage()
			} else {
				// Initialize the default help flag for the found command and display its help information.
				cmd.InitDefaultHelpFlag()
				_ = cmd.Help()
			}
		},
	}
}

// addHelpCommandFlag adds a help flag to the provided flag set.
// The flag provides contextual help based on the usage of the command.
func addHelpCommandFlag(usage string, fs *pflag.FlagSet) {
	fs.BoolP(
		flagHelp,
		flagHelpShorthand,
		false,
		// Dynamically generate the help description based on the usage.
		fmt.Sprintf("Display help for the %s command.", color.GreenString(strings.Split(usage, " ")[0])),
	)
}
