package log

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	flagLogLevel             = "log.level"
	flagLogDisableCaller     = "log.disable-caller"
	flagLogDisableStacktrace = "log.disable-stacktrace"
	flagLogFormat            = "log.format"
	flagLogEnableColor       = "log.enable-color"
	flagLogOutputPaths       = "log.output-paths"
	flagLogErrorOutputPaths  = "log.error-output-paths"
	flagLogDevelopment       = "log.development"
	flagLogName              = "log.name"

	consoleFormat = "console"
	jsonFormat    = "json"
)

type Options struct {
	Level             string   `json:"level"              mapstructure:"level"`
	Format            string   `json:"format"             mapstructure:"format"`
	DisableCaller     bool     `json:"disable-caller"     mapstructure:"disable-caller"`
	DisableStacktrace bool     `json:"disable-stacktrace" mapstructure:"disable-stacktrace"`
	EnableColor       bool     `json:"enable-color"       mapstructure:"enable-color"`
	Development       bool     `json:"development"        mapstructure:"development"`
	Name              string   `json:"name"               mapstructure:"name"`
	OutputPaths       []string `json:"output-paths"       mapstructure:"output-paths"`
	ErrorOutputPaths  []string `json:"error-output-paths" mapstructure:"error-output-paths"`
}

// NewOptions creates a new instance of Options with default values.
func NewOptions() *Options {
	return &Options{
		Level:             zapcore.InfoLevel.String(),
		Format:            consoleFormat,
		DisableCaller:     false,
		DisableStacktrace: false,
		EnableColor:       false,
		Development:       false,
		Name:              "default",
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
	}
}

// Validate checks the validity of the Options configuration. It ensures
// that the log level and format are set correctly.
func (o *Options) Validate() []error {
	var errs []error

	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(o.Level)); err != nil {
		errs = append(errs, err)
	}

	format := strings.ToLower(o.Format)
	if format != consoleFormat && format != jsonFormat {
		errs = append(errs, fmt.Errorf("invalid log format: %q", o.Format))
	}

	return errs
}

// String returns a JSON string representation of the Options. It is useful
// for logging and debugging purposes.
func (o *Options) String() string {
	res, _ := json.Marshal(o)
	return string(res)
}

// AddFlags adds logging-related command line flags to the provided FlagSet.
// This allows the configuration of the logging system via command line parameters.
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Level, flagLogLevel, o.Level, "Define the minimum log output level.")
	fs.BoolVar(&o.DisableCaller, flagLogDisableCaller, o.DisableCaller, "Disable output of caller information in logs.")
	fs.BoolVar(&o.DisableStacktrace, flagLogDisableStacktrace, o.DisableStacktrace, "Disable the logging of stack trace at or above the panic level.")
	fs.StringVar(&o.Format, flagLogFormat, o.Format, "Define the log output format, either console or json.")
	fs.BoolVar(&o.EnableColor, flagLogEnableColor, o.EnableColor, "Enable colored output for console format logs.")
	fs.StringSliceVar(&o.OutputPaths, flagLogOutputPaths, o.OutputPaths, "Define the log output paths.")
	fs.StringSliceVar(&o.ErrorOutputPaths, flagLogErrorOutputPaths, o.ErrorOutputPaths, "Define the error log output paths.")
	fs.BoolVar(&o.Development, flagLogDevelopment, o.Development, "Enable development mode for more verbose logging output.")
	fs.StringVar(&o.Name, flagLogName, o.Name, "Set a name for the logger.")
}

// Build constructs and configures a global zap logger based on the Options settings.
// It sets up the logger with the specified level, format, output paths, and other configurations.
func (o *Options) Build() error {
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(o.Level)); err != nil {
		return fmt.Errorf("invalid log level: %v", err)
	}

	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	if o.Format == consoleFormat && !o.EnableColor {
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	}

	zapConfig := zap.Config{
		Level:             zap.NewAtomicLevelAt(zapLevel),
		Development:       o.Development,
		DisableCaller:     o.DisableCaller,
		DisableStacktrace: o.DisableStacktrace,
		Encoding:          o.Format,
		EncoderConfig:     encoderConfig,
		OutputPaths:       o.OutputPaths,
		ErrorOutputPaths:  o.ErrorOutputPaths,
	}

	logger, err := zapConfig.Build()
	if err != nil {
		return fmt.Errorf("failed to build logger: %v", err)
	}

	zap.RedirectStdLog(logger.Named(o.Name))
	zap.ReplaceGlobals(logger)

	return nil
}
