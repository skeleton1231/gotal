// Copyright 2023 Talhuang<talhuang1231@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package app

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// Command struct represents an individual CLI command.
type Command struct {
	usage    string         // A brief string about how the command is used.
	desc     string         // A description of the command.
	options  CliOptions     // Options associated with the command.
	commands []*Command     // Subcommands under this command.
	runFunc  RunCommandFunc // Function to run when the command is executed.
}

// CommandOption represents a function that configures a Command.
type CommandOption func(*Command)

// WithCommandOptions returns a CommandOption that sets the options for a command.
func WithCommandOptions(opt CliOptions) CommandOption {
	return func(c *Command) {
		c.options = opt
	}
}

// RunCommandFunc defines the function signature for commands that run.
type RunCommandFunc func(args []string) error

// WithCommandRunFunc returns a CommandOption that sets the run function of a command.
func WithCommandRunFunc(run RunCommandFunc) CommandOption {
	return func(c *Command) {
		c.runFunc = run
	}
}

// NewCommand initializes and returns a new Command.
func NewCommand(usage string, desc string, opts ...CommandOption) *Command {
	c := &Command{
		usage: usage,
		desc:  desc,
	}

	// Apply all provided command options.
	for _, o := range opts {
		o(c)
	}

	return c
}

// AddCommand adds a subcommand to the current command.
func (c *Command) AddCommand(cmd *Command) {
	c.commands = append(c.commands, cmd)
}

// AddCommands adds multiple subcommands to the current command.
func (c *Command) AddCommands(cmds ...*Command) {
	c.commands = append(c.commands, cmds...)
}

// cobraCommand converts the Command to a *cobra.Command.
func (c *Command) cobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   c.usage,
		Short: c.desc,
	}
	cmd.SetOutput(os.Stdout)
	cmd.Flags().SortFlags = false
	if len(c.commands) > 0 {
		for _, command := range c.commands {
			cmd.AddCommand(command.cobraCommand())
		}
	}
	if c.runFunc != nil {
		cmd.Run = c.runCommand
	}
	if c.options != nil {
		for _, f := range c.options.Flags().FlagSets {
			cmd.Flags().AddFlagSet(f)
		}
	}
	addHelpCommandFlag(c.usage, cmd.Flags())

	return cmd
}

// runCommand is the function to run when the cobra command is executed.
func (c *Command) runCommand(cmd *cobra.Command, args []string) {
	if c.runFunc != nil {
		if err := c.runFunc(args); err != nil {
			// Print the error and exit the program with an error code.
			fmt.Printf("%v %v\n", color.RedString("Error:"), err)
			os.Exit(1)
		}
	}
}

// App functions were not defined in the given code, but are used for adding commands to an App instance.

// AddCommand adds a command to the App.
func (a *App) AddCommand(cmd *Command) {
	a.commands = append(a.commands, cmd)
}

// AddCommands adds multiple commands to the App.
func (a *App) AddCommands(cmds ...*Command) {
	a.commands = append(a.commands, cmds...)
}

// FormatBaseName cleans up and formats the basename, especially for windows OS.
func FormatBaseName(basename string) string {
	if runtime.GOOS == "windows" {
		basename = strings.ToLower(basename)
		basename = strings.TrimSuffix(basename, ".exe")
	}

	return basename
}
