// Copyright 2023 Talhuang<talhuang1231@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package logger

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	gormlogger "gorm.io/gorm/logger"
)

// Define colors.
const (
	Reset       = "\033[0m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	Cyan        = "\033[36m"
	White       = "\033[37m"
	BlueBold    = "\033[34;1m"
	MagentaBold = "\033[35;1m"
	RedBold     = "\033[31;1m"
	YellowBold  = "\033[33;1m"
)

// Define gorm log levels.
const (
	Silent gormlogger.LogLevel = iota + 1
	Error
	Warn
	Info
)

// Writer log writer interface.
type Writer interface {
	Printf(string, ...interface{})
}

// Config defines a gorm logger configuration.
type Config struct {
	SlowThreshold time.Duration
	Colorful      bool
	LogLevel      gormlogger.LogLevel
}

// New creates a gorm logger instance with logrus as the underlying logger.
func New(level int) gormlogger.Interface {
	// Setup string templates for log messages, these are kept
	// as they are used to format the message before sending to logrus.
	// Note that coloring will be handled by logrus' formatter.

	config := Config{
		SlowThreshold: 200 * time.Millisecond,
		Colorful:      true, // assuming you want colors with logrus formatter
		LogLevel:      gormlogger.LogLevel(level),
	}

	// Configure logrus' formatter here if needed.

	return &logger{
		Config: config,
		// Color strings removed from the format strings since logrus handles it.
		infoStr:      "[info] ",
		warnStr:      "[warn] ",
		errStr:       "[error] ",
		traceStr:     "[%.3fms] [rows:%v] %s",
		traceWarnStr: "%s[%.3fms] [rows:%v] %s",
		traceErrStr:  "%s[%.3fms] [rows:%v] %s",
	}
}

type logger struct {
	Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

// LogMode sets the log level for the logger.
func (l *logger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info logs info level messages using logrus.
func (l logger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= Info {
		logrus.Infof(l.infoStr+msg, append([]interface{}{fileWithLineNum()}, data...)...)
	}
}

// Warn logs warning level messages using logrus.
func (l logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= Warn {
		logrus.Warnf(l.warnStr+msg, append([]interface{}{fileWithLineNum()}, data...)...)
	}
}

// Error logs error level messages using logrus.
func (l logger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= Error {
		logrus.Errorf(l.errStr+msg, append([]interface{}{fileWithLineNum()}, data...)...)
	}
}

// Trace logs SQL queries using logrus.
func (l logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	// Early return if logging is disabled.
	if l.LogLevel <= Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	formattedElapsed := float64(elapsed.Nanoseconds()) / 1e6

	if err != nil && l.LogLevel >= Error {
		// Log as error if there is one.
		if rows == -1 {
			logrus.Errorf(l.traceErrStr, fileWithLineNum(), err, formattedElapsed, "-", sql)
		} else {
			logrus.Errorf(l.traceErrStr, fileWithLineNum(), err, formattedElapsed, rows, sql)
		}
	} else if elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= Warn {
		// Log as warning if the query is slow.
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			logrus.Warnf(l.traceWarnStr, fileWithLineNum(), slowLog, formattedElapsed, "-", sql)
		} else {
			logrus.Warnf(l.traceWarnStr, fileWithLineNum(), slowLog, formattedElapsed, rows, sql)
		}
	} else if l.LogLevel >= Info {
		// Otherwise, log as info.
		if rows == -1 {
			logrus.Infof(l.traceStr, fileWithLineNum(), formattedElapsed, "-", sql)
		} else {
			logrus.Infof(l.traceStr, fileWithLineNum(), formattedElapsed, rows, sql)
		}
	}
}

// fileWithLineNum returns the file and line number of the calling function.
func fileWithLineNum() string {
	for i := 4; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok && !strings.HasSuffix(file, "_test.go") {
			dir, f := filepath.Split(file)
			return filepath.Join(filepath.Base(dir), f) + ":" + strconv.FormatInt(int64(line), 10)
		}
	}
	return ""
}
