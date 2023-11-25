package log

import (
	"context"
	"fmt"
	"log"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type InfoLogger interface {
	Info(msg string, fields ...Field)
	Infof(format string, v ...interface{})
	Infokv(msg string, kv ...interface{})
	Enabled() bool
}

type Logger interface {
	InfoLogger
	Debug(msg string, fields ...Field)
	Debugf(format string, v ...interface{})
	Debugkv(msg string, kv ...interface{})
	Warn(msg string, fields ...Field)
	Warnf(format string, v ...interface{})
	Warnkv(msg string, kv ...interface{})
	Error(msg string, fields ...Field)
	Errorf(format string, v ...interface{})
	Errorkv(msg string, kv ...interface{})
	Panic(msg string, fields ...Field)
	Panicf(format string, v ...interface{})
	Panickv(msg string, kv ...interface{})
	Fatal(msg string, fields ...Field)
	Fatalf(format string, v ...interface{})
	Fatalkv(msg string, kv ...interface{})

	LoggerForLevel(level zapcore.Level) InfoLogger
	Write(p []byte) (n int, err error)

	WithValues(kv ...interface{}) Logger

	WithName(name string) Logger

	WithContext(ctx context.Context) context.Context

	Sync()
}

type infoLogger struct {
	level zapcore.Level
	log   *zap.Logger
}

func (l *infoLogger) Enabled() bool { return true }
func (l *infoLogger) Info(msg string, fields ...Field) {
	if checkedEntry := l.log.Check(l.level, msg); checkedEntry != nil {
		checkedEntry.Write(fields...)
	}
}

func (l *infoLogger) Infof(format string, args ...interface{}) {
	if checkedEntry := l.log.Check(l.level, fmt.Sprintf(format, args...)); checkedEntry != nil {
		checkedEntry.Write()
	}
}

func (l *infoLogger) Infokv(msg string, kv ...interface{}) {
	if checkedEntry := l.log.Check(l.level, msg); checkedEntry != nil {
		checkedEntry.Write(processFields(l.log, kv)...)
	}
}

func processFields(logger *zap.Logger, args []interface{}, additionalFields ...zap.Field) []zap.Field {
	if len(args)%2 != 0 {
		logger.DPanic("expected even number of arguments for key-value pairs", zap.Int("received_args", len(args)))
	}

	fields := make([]zap.Field, 0, len(args)/2+len(additionalFields))
	for i := 0; i < len(args); i += 2 {
		if i+1 >= len(args) {
			logger.DPanic("missing value for the key", zap.Any("key", args[i]))
			break
		}

		key, val := args[i], args[i+1]
		keyStr, ok := key.(string)
		if !ok {
			logger.DPanic("key is not a string", zap.Any("received_key", key))
			continue
		}

		fields = append(fields, zap.Any(keyStr, val))
	}

	return append(fields, additionalFields...)
}

var (
	logger = New(NewOptions())
	mu     sync.Mutex
)

func Init(opts *Options) {
	mu.Lock()
	defer mu.Unlock()
	logger = New(opts)
}

// New create logger by opts which can custmoized by command arguments.
func New(opts *Options) *zapLogger {
	if opts == nil {
		opts = NewOptions()
	}

	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(opts.Level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}
	encodeLevel := zapcore.CapitalLevelEncoder
	// when output to local path, with color is forbidden
	if opts.Format == consoleFormat && opts.EnableColor {
		encodeLevel = zapcore.CapitalColorLevelEncoder
	}

	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "timestamp",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    encodeLevel,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: milliSecondsDurationEncoder,
	}

	loggerConfig := &zap.Config{
		Level:             zap.NewAtomicLevelAt(zapLevel),
		Development:       opts.Development,
		DisableCaller:     opts.DisableCaller,
		DisableStacktrace: opts.DisableStacktrace,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         opts.Format,
		EncoderConfig:    encoderConfig,
		OutputPaths:      opts.OutputPaths,
		ErrorOutputPaths: opts.ErrorOutputPaths,
	}

	var err error
	l, err := loggerConfig.Build(zap.AddStacktrace(zapcore.PanicLevel), zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
	logger := &zapLogger{
		zapLogger: l.Named(opts.Name),
		infoLogger: infoLogger{
			log:   l,
			level: zap.InfoLevel,
		},
	}
	zap.RedirectStdLog(l)

	return logger
}

// SugaredLogger returns a *zap.SugaredLogger instance, offering a more flexible logging approach.
// This logger supports fast and loosely-typed logging, suitable for most use cases.
func LoggerSugared() *zap.SugaredLogger {
	return logger.zapLogger.Sugar()
}

// ErrLogger creates and returns a standard library *log.Logger instance that logs at the ErrorLevel in zap.
// If the global logger is not initialized, it safely returns nil.
func LoggerErr() *log.Logger {
	if logger == nil {
		return nil
	}
	if l, err := zap.NewStdLogAt(logger.zapLogger, zapcore.ErrorLevel); err == nil {
		return l
	}

	return nil
}

// InfoLogger creates and returns a standard library *log.Logger instance that logs at the InfoLevel in zap.
// If the global logger is not initialized, it returns nil.
func LoggerInfo() *log.Logger {
	if logger == nil {
		return nil
	}
	if l, err := zap.NewStdLogAt(logger.zapLogger, zapcore.InfoLevel); err == nil {
		return l
	}

	return nil
}

type zapLogger struct {
	zapLogger *zap.Logger
	infoLogger
}

func (l *zapLogger) Debug(msg string, fields ...Field) {
	l.zapLogger.Debug(msg, fields...)
}

func (l *zapLogger) Debugf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Debugf(format, v...)
}

func (l *zapLogger) Debugkv(msg string, kv ...interface{}) {
	l.zapLogger.Sugar().Debugw(msg, kv...)
}

func (l *zapLogger) Info(msg string, fields ...Field) {
	l.zapLogger.Info(msg, fields...)
}

func (l *zapLogger) Infof(format string, v ...interface{}) {
	l.zapLogger.Sugar().Infof(format, v...)
}

func (l *zapLogger) Infokv(msg string, kv ...interface{}) {
	l.zapLogger.Sugar().Infow(msg, kv...)
}

func (l *zapLogger) Warn(msg string, fields ...Field) {
	l.zapLogger.Warn(msg, fields...)
}

func (l *zapLogger) Warnf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Warnf(format, v...)
}

func (l *zapLogger) Warnkv(msg string, kv ...interface{}) {
	l.zapLogger.Sugar().Warnw(msg, kv...)
}

func (l *zapLogger) Error(msg string, fields ...Field) {
	l.zapLogger.Error(msg, fields...)
}

func (l *zapLogger) Errorf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Errorf(format, v...)
}

func (l *zapLogger) Errorkv(msg string, kv ...interface{}) {
	l.zapLogger.Sugar().Errorw(msg, kv...)
}

func (l *zapLogger) Panic(msg string, fields ...Field) {
	l.zapLogger.Panic(msg, fields...)
}

func (l *zapLogger) Panicf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Panicf(format, v...)
}

func (l *zapLogger) Panickv(msg string, kv ...interface{}) {
	l.zapLogger.Sugar().Panicw(msg, kv...)
}

func (l *zapLogger) Fatal(msg string, fields ...Field) {
	l.zapLogger.Fatal(msg, fields...)
}

func (l *zapLogger) Fatalf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Fatalf(format, v...)
}

func (l *zapLogger) Fatalkv(msg string, kv ...interface{}) {
	l.zapLogger.Sugar().Fatalw(msg, kv...)
}

func (l *zapLogger) LoggerForLevel(level zapcore.Level) InfoLogger {
	// 根据 zapcore.Level 返回相应级别的 logger
	return &infoLogger{level: level, log: l.zapLogger}
}

func (l *zapLogger) Write(p []byte) (n int, err error) {
	l.zapLogger.Info(string(p))
	return len(p), nil
}

func (l *zapLogger) Sync() {
	_ = l.zapLogger.Sync()
}

func (l *zapLogger) WithValues(kv ...interface{}) Logger {
	newLogger := l.zapLogger.With(processFields(l.zapLogger, kv)...)
	return &zapLogger{zapLogger: newLogger}
}

func (l *zapLogger) WithName(name string) Logger {
	newLogger := l.zapLogger.Named(name)
	return &zapLogger{zapLogger: newLogger}
}

// Record extracts additional context (like request ID, username, etc.) from the provided context.Context
// and returns a new logger instance with this context.
func (l *zapLogger) Record(ctx context.Context) *zapLogger {
	lg := l.clone()

	if requestID := ctx.Value(KeyRequestID); requestID != nil {
		lg.zapLogger = lg.zapLogger.With(zap.Any(KeyRequestID, requestID))
	}
	if username := ctx.Value(KeyUsername); username != nil {
		lg.zapLogger = lg.zapLogger.With(zap.Any(KeyUsername, username))
	}

	return lg
}

func (l *zapLogger) clone() *zapLogger {
	copy := *l

	return &copy
}

func NewLogger(l *zap.Logger) Logger {
	return &zapLogger{
		zapLogger: l,
		infoLogger: infoLogger{
			log:   l,
			level: zap.InfoLevel,
		},
	}
}

// LoggerForLevel returns a logger for the specified logging level.
func LoggerForLevel(level Level) InfoLogger { return logger.LoggerForLevel(level) }

func Debug(msg string, fields ...Field) {
	logger.zapLogger.Debug(msg, fields...)
}

func Debugf(format string, v ...interface{}) {
	logger.zapLogger.Sugar().Debugf(format, v...)
}

func Debugw(msg string, kv ...interface{}) {
	logger.zapLogger.Sugar().Debugw(msg, kv...)
}

func Info(msg string, fields ...Field) {
	logger.zapLogger.Info(msg, fields...)
}

func Infof(format string, v ...interface{}) {
	logger.zapLogger.Sugar().Infof(format, v...)
}

func Infow(msg string, kv ...interface{}) {
	logger.zapLogger.Sugar().Infow(msg, kv...)
}

func Warn(msg string, fields ...Field) {
	logger.zapLogger.Warn(msg, fields...)
}

func Warnf(format string, v ...interface{}) {
	logger.zapLogger.Sugar().Warnf(format, v...)
}

func Warnw(msg string, kv ...interface{}) {
	logger.zapLogger.Sugar().Warnw(msg, kv...)
}

func Error(msg string, fields ...Field) {
	logger.zapLogger.Error(msg, fields...)
}

func Errorf(format string, v ...interface{}) {
	logger.zapLogger.Sugar().Errorf(format, v...)
}

func Errorw(msg string, kv ...interface{}) {
	logger.zapLogger.Sugar().Errorw(msg, kv...)
}

func Panic(msg string, fields ...Field) {
	logger.zapLogger.Panic(msg, fields...)
}

func Panicf(format string, v ...interface{}) {
	logger.zapLogger.Sugar().Panicf(format, v...)
}

func Panicw(msg string, kv ...interface{}) {
	logger.zapLogger.Sugar().Panicw(msg, kv...)
}

func Fatal(msg string, fields ...Field) {
	logger.zapLogger.Fatal(msg, fields...)
}

func Fatalf(format string, v ...interface{}) {
	logger.zapLogger.Sugar().Fatalf(format, v...)
}

func Fatalw(msg string, kv ...interface{}) {
	logger.zapLogger.Sugar().Fatalw(msg, kv...)
}

func Record(ctx context.Context) *zapLogger {
	return logger.Record(ctx)
}

func WithValues(kv ...interface{}) Logger { return logger.WithValues(kv...) }

func WithName(s string) Logger { return logger.WithName(s) }

func Flush() { logger.Sync() }
