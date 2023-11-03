// Copyright 2023 Talhuang<talhuang1231@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package app

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"

	"github.com/skeleton1231/gotal/pkg/util/flag"
	globalflag "github.com/skeleton1231/gotal/pkg/util/flag"
	sections "github.com/skeleton1231/gotal/pkg/util/flag"
	"github.com/skeleton1231/gotal/pkg/util/term"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	progressMessage = color.GreenString("==>")

	usageTemplate = fmt.Sprintf(`%s{{if .Runnable}}
  %s{{end}}{{if .HasAvailableSubCommands}}
  %s{{end}}{{if gt (len .Aliases) 0}}

%s
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

%s
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

%s{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  %s {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

%s
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

%s
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

%s{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "%s --help" for more information about a command.{{end}}
`,
		color.CyanString("Usage:"),
		color.GreenString("{{.UseLine}}"),
		color.GreenString("{{.CommandPath}} [command]"),
		color.CyanString("Aliases:"),
		color.CyanString("Examples:"),
		color.CyanString("Available Commands:"),
		color.GreenString("{{rpad .Name .NamePadding }}"),
		color.CyanString("Flags:"),
		color.CyanString("Global Flags:"),
		color.CyanString("Additional help topics:"),
		color.GreenString("{{.CommandPath}} [command]"),
	)
)

// App is the central structure representing the CLI application.
type App struct {
	basename    string
	name        string
	description string
	options     CliOptions
	runFunc     RunFunc
	commands    []*Command
	args        cobra.PositionalArgs
	cmd         *cobra.Command
}

// Option represents a function that can modify an App.
type Option func(*App)

// WithOptions allows CLI options to be passed into the App.
func WithOptions(opt CliOptions) Option {
	return func(a *App) {
		a.options = opt
	}
}

// RunFunc is a function callback that is executed when the app starts.
type RunFunc func(basename string) error

// WithRunFunc sets the run function callback for the app.
func WithRunFunc(run RunFunc) Option {
	return func(a *App) {
		a.runFunc = run
	}
}

// WithDescription sets a description for the App.
func WithDescription(desc string) Option {
	return func(a *App) {
		a.description = desc
	}
}

// WithValidArgs sets validation for command arguments.
func WithValidArgs(args cobra.PositionalArgs) Option {
	return func(a *App) {
		a.args = args
	}
}

// WithDefaultValidArgs sets default validation for command arguments.
func WithDefaultValidArgs() Option {
	return func(a *App) {
		a.args = func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}
			return nil
		}
	}
}

// NewApp initializes a new App instance.
func NewApp(name string, basename string, opts ...Option) *App {

	a := &App{
		name:     name,
		basename: basename,
	}
	// Apply options to the App instance.
	for _, o := range opts {
		o(a)
	}
	// Build the underlying cobra command.
	a.buildCommand()
	return a
}

// buildCommand constructs the main cobra command for the App.
func (a *App) buildCommand() {
	cmd := cobra.Command{
		Use:           FormatBaseName(a.basename),
		Short:         a.name,
		Long:          a.description,
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          a.args,
	}
	// Basic settings for command output.
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	// Sort the flags for better visualization.
	cmd.Flags().SortFlags = true
	flag.InitFlags(cmd.Flags())

	// Add sub-commands to the main command.
	if len(a.commands) > 0 {
		for _, command := range a.commands {
			cmd.AddCommand(command.cobraCommand())
		}
		// Set a custom help command.
		cmd.SetHelpCommand(helpCommand(FormatBaseName(a.basename)))
	}
	// Set the main run function.
	if a.runFunc != nil {
		cmd.RunE = a.runCommand
	}

	// Add flags to the command.
	var namedFlagSets sections.NamedFlagSets
	if a.options != nil {
		namedFlagSets = a.options.Flags()
		fs := cmd.Flags()
		for _, f := range namedFlagSets.FlagSets {
			fs.AddFlagSet(f)
		}
	}

	// Init Global Config
	addConfigFlag(a.basename, namedFlagSets.FlagSet("global"))
	globalflag.AddGlobalFlags(namedFlagSets.FlagSet("global"), cmd.Name())
	cmd.Flags().AddFlagSet(namedFlagSets.FlagSet("global"))

	// Apply custom usage and help templates.
	addCmdTemplate(&cmd, namedFlagSets)
	// Assign the main cobra command to the App instance.
	a.cmd = &cmd
}

// Run starts the App.
func (a *App) Run() {
	if err := a.cmd.Execute(); err != nil {
		fmt.Printf("%v %v\n", color.RedString("Error:"), err)
		os.Exit(1)
	}
}

// Command returns the main cobra command for the App.
func (a *App) Command() *cobra.Command {
	return a.cmd
}

// runCommand executes when the main command is run.
func (a *App) runCommand(cmd *cobra.Command, args []string) error {
	printWorkingDir()
	flag.PrintFlags(cmd.Flags())

	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	if err := viper.Unmarshal(a.options); err != nil {
		return err
	}

	logrus.Infof("%v Starting %s ...", progressMessage, a.name)
	// Apply option rules.
	if a.options != nil {
		if err := a.applyOptionRules(); err != nil {
			return err
		}
	}
	// Run the main application function.
	if a.runFunc != nil {
		return a.runFunc(a.basename)
	}
	return nil
}

// applyOptionRules applies any additional rules or configurations for options.
func (a *App) applyOptionRules() error {
	if completeableOptions, ok := a.options.(CompleteableOptions); ok {
		if err := completeableOptions.Complete(); err != nil {
			return err
		}
	}

	if printableOptions, ok := a.options.(PrintableOptions); ok {
		logrus.Infof("%v Config: `%s`", progressMessage, printableOptions.String())
	}
	return nil
}

// printWorkingDir prints the current working directory to the console.
func printWorkingDir() {
	wd, _ := os.Getwd()
	logrus.Infof("%v WorkingDir: %s", progressMessage, wd)
}

// addCmdTemplate customizes the usage and help templates for the command.
func addCmdTemplate(cmd *cobra.Command, namedFlagSets sections.NamedFlagSets) {
	usageFmt := "Usage:\n  %s\n"
	cols, _, _ := term.TerminalSize(cmd.OutOrStdout())
	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		fmt.Fprintf(cmd.OutOrStderr(), usageFmt, cmd.UseLine())
		sections.PrintSections(cmd.OutOrStderr(), namedFlagSets, cols)

		return nil
	})
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
		sections.PrintSections(cmd.OutOrStdout(), namedFlagSets, cols)

	})
}
