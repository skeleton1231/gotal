package log

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type InfoLogger interface {
	Info(msg string, fields ...zapcore.Field)
	Infof(format string, v ...interface{})
	Infokv(msg string, kv ...interface{})
	Enabled() bool
}

type Logger interface {
	InfoLogger
	Debug(msg string, fields ...zapcore.Field)
	Debugf(format string, v ...interface{})
	Debugkv(msg string, kv ...interface{})
	Warn(msg string, fields ...zapcore.Field)
	Warnf(format string, v ...interface{})
	Warnkv(msg string, kv ...interface{})
	Error(msg string, fields ...zapcore.Field)
	Errorf(format string, v ...interface{})
	Errorkv(msg string, kv ...interface{})
	Panic(msg string, fields ...zapcore.Field)
	Panicf(format string, v ...interface{})
	Panickv(msg string, kv ...interface{})
	Fatal(msg string, fields ...zapcore.Field)
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
func (l *infoLogger) Info(msg string, fields ...zapcore.Field) {
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

type zapLogger struct {
	zapLogger *zap.Logger
	infoLogger
}

func (l *zapLogger) Debug(msg string, fields ...zapcore.Field) {
	l.zapLogger.Debug(msg, fields...)
}

func (l *zapLogger) Debugf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Debugf(format, v...)
}

func (l *zapLogger) Debugkv(msg string, kv ...interface{}) {
	l.zapLogger.Sugar().Debugw(msg, kv...)
}

func (l *zapLogger) Info(msg string, fields ...zapcore.Field) {
	l.zapLogger.Info(msg, fields...)
}

func (l *zapLogger) Infof(format string, v ...interface{}) {
	l.zapLogger.Sugar().Infof(format, v...)
}

func (l *zapLogger) Infokv(msg string, kv ...interface{}) {
	l.zapLogger.Sugar().Infow(msg, kv...)
}

func (l *zapLogger) Warn(msg string, fields ...zapcore.Field) {
	l.zapLogger.Warn(msg, fields...)
}

func (l *zapLogger) Warnf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Warnf(format, v...)
}

func (l *zapLogger) Warnkv(msg string, kv ...interface{}) {
	l.zapLogger.Sugar().Warnw(msg, kv...)
}

func (l *zapLogger) Error(msg string, fields ...zapcore.Field) {
	l.zapLogger.Error(msg, fields...)
}

func (l *zapLogger) Errorf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Errorf(format, v...)
}

func (l *zapLogger) Errorkv(msg string, kv ...interface{}) {
	l.zapLogger.Sugar().Errorw(msg, kv...)
}

func (l *zapLogger) Panic(msg string, fields ...zapcore.Field) {
	l.zapLogger.Panic(msg, fields...)
}

func (l *zapLogger) Panicf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Panicf(format, v...)
}

func (l *zapLogger) Panickv(msg string, kv ...interface{}) {
	l.zapLogger.Sugar().Panicw(msg, kv...)
}

func (l *zapLogger) Fatal(msg string, fields ...zapcore.Field) {
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

func NewLogger(l *zap.Logger) Logger {
	return &zapLogger{
		zapLogger: l,
		infoLogger: infoLogger{
			log:   l,
			level: zap.InfoLevel,
		},
	}
}
