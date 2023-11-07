// Copyright 2023 Talhuang<talhuang1231@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package middleware provides middleware utilities for the Gin web framework.
package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mattn/go-isatty"
	"github.com/sirupsen/logrus"
)

// defaultLogFormatter defines the default format for logging HTTP requests.
var defaultLogFormatter = func(param gin.LogFormatterParams) string {
	// Returns the formatted log string.
	return fmt.Sprintf("%3d - [%s] \"%v %s %s\" %s",
		param.StatusCode,
		param.ClientIP,
		param.Latency,
		param.Method,
		param.Path,
		param.ErrorMessage,
	)
}

// Logger returns a middleware handler that logs HTTP requests.
func Logger() gin.HandlerFunc {
	// Returns a Logger with the default Gin logger configuration.
	return LoggerWithConfig(GetLoggerConfig(nil, nil, nil))
}

// LoggerWithConfig takes a gin.LoggerConfig and returns a logging middleware.
func LoggerWithConfig(conf gin.LoggerConfig) gin.HandlerFunc {
	// Chooses which formatter to use for logging.
	formatter := conf.Formatter
	if formatter == nil {
		formatter = defaultLogFormatter
	}

	// Determines the output destination for the logs.
	out := conf.Output
	if out == nil {
		out = gin.DefaultWriter
	}

	// Checks if output should be colored.
	isTerm := shouldUseColor(out)

	// Enables colored output if supported.
	if isTerm {
		gin.ForceConsoleColor()
	}

	// Creates a map of paths that should not be logged.
	skip := skippedPaths(conf.SkipPaths)

	// The middleware function that does the logging.
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		// Processes the request and calls the next handler.
		c.Next()

		// Checks if the current path is not in the skip list.
		if _, ok := skip[path]; !ok {
			// Gathers log parameters from the request.
			param, contextFields := createLogFormatterParams(c, start, path)

			// Marshal the context fields to a JSON string.
			contextJSON, err := json.Marshal(contextFields)
			if err != nil {
				// Handle error, perhaps log that the context fields could not be marshaled.
				logrus.WithError(err).Error("Could not marshal context fields")
			}

			logrus.WithFields(logrus.Fields{
				"client_ip":     param.ClientIP,
				"method":        param.Method,
				"status_code":   param.StatusCode,
				"latency":       param.Latency,
				"user_agent":    param.Request.UserAgent(),
				"error_message": param.ErrorMessage,
				"body_size":     param.BodySize,
				"path":          param.Path,
				"context":       string(contextJSON), // Add the JSON string as the "context" field
			}).Info("Request handled")

		}
	}
}

// shouldUseColor determines if the log output should be colored based on the output writer.
func shouldUseColor(out io.Writer) bool {
	// Checks if the output is a terminal and if it supports color.
	if w, ok := out.(*os.File); !ok || os.Getenv("TERM") == "dumb" || (!isatty.IsTerminal(w.Fd()) && !isatty.IsCygwinTerminal(w.Fd())) {
		return false
	}
	return true
}

// skippedPaths creates a map from a list of paths to be skipped by the logger.
func skippedPaths(notlogged []string) map[string]struct{} {
	skip := make(map[string]struct{})
	for _, path := range notlogged {
		// Adds the path to the map to indicate it should be skipped.
		skip[path] = struct{}{}
	}
	return skip
}

// createLogFormatterParams gathers information from the context to create log formatter parameters.
func createLogFormatterParams(c *gin.Context, start time.Time, path string) (gin.LogFormatterParams, logrus.Fields) {
	// Initializes the parameters structure with data from the request.
	param := gin.LogFormatterParams{
		Request: c.Request,
		Keys:    c.Keys,
	}

	// Completes the parameter structure with additional timing and request info.
	param.TimeStamp = time.Now()
	param.Latency = param.TimeStamp.Sub(start)
	param.ClientIP = c.ClientIP()
	param.Method = c.Request.Method
	param.StatusCode = c.Writer.Status()
	param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
	param.BodySize = c.Writer.Size()
	param.Path = path + getRawQuery(c)

	// Gather Logrus fields from the context.
	logrusFields := logrus.Fields{}
	if requestID, exists := c.Get("requestID"); exists {
		logrusFields["requestID"] = requestID
	}
	if username, exists := c.Get("username"); exists {
		logrusFields["username"] = username
	}

	return param, logrusFields
}

// getRawQuery appends the raw query string to the path if present.
func getRawQuery(c *gin.Context) string {
	raw := c.Request.URL.RawQuery
	if raw != "" {
		return "?" + raw
	}
	return ""
}

// // createLogEntry generates a structured logrus entry from the log formatter parameters.
// func createLogEntry(param gin.LogFormatterParams) *logrus.Entry {
// 	// Maps the log formatter parameters to a structured log entry.
// 	return logrus.WithFields(logrus.Fields{
// 		"client_ip":     param.ClientIP,
// 		"method":        param.Method,
// 		"status_code":   param.StatusCode,
// 		"latency":       param.Latency,
// 		"user_agent":    param.Request.UserAgent(),
// 		"error_message": param.ErrorMessage,
// 		"body_size":     param.BodySize,
// 		"path":          param.Path,
// 	})
// }
