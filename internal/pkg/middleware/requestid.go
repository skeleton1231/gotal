// Copyright 2023 Talhuang<talhuang1231@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package middleware

import (
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

const (
	// XRequestIDKey is the header field used to store the request ID.
	XRequestIDKey = "X-Request-ID"
)

// RequestID is a middleware that injects a request ID into each request.
// If the incoming request already has an 'X-Request-ID' header, it will use that one,
// otherwise, it generates a new UUIDv4 and sets it in both the request header and context.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the request ID from the incoming request header.
		rid := c.GetHeader(XRequestIDKey)

		// If there is no request ID present, generate a new UUIDv4.
		if rid == "" {
			rid = uuid.NewV4().String()
			// Set the generated UUID in the request header.
			c.Request.Header.Set(XRequestIDKey, rid)
			// Store the request ID in the Gin context for later use.
			c.Set(XRequestIDKey, rid)
		}

		// Set the request ID in the response header.
		c.Writer.Header().Set(XRequestIDKey, rid)
		// Proceed with the next middleware or handler.
		c.Next()
	}
}

// GetLoggerConfig creates and returns a gin.LoggerConfig configured to write logs to a specified io.Writer.
// The logs will be formatted using the provided gin.LogFormatter.
// If no formatter is provided, it defaults to a custom log formatter that includes the request ID.
// gin.DefaultWriter is os.Stdout by default, but can be set to any io.Writer.
func GetLoggerConfig(formatter gin.LogFormatter, output io.Writer, skipPaths []string) gin.LoggerConfig {
	// If no custom formatter is provided, use the default one that includes the request ID.
	if formatter == nil {
		formatter = GetDefaultLogFormatterWithRequestID()
	}

	// Return a LoggerConfig with the provided formatter, output destination, and paths to skip logging.
	return gin.LoggerConfig{
		Formatter: formatter,
		Output:    output,
		SkipPaths: skipPaths,
	}
}

// GetDefaultLogFormatterWithRequestID creates a log formatter function that includes the request ID in the log message.
func GetDefaultLogFormatterWithRequestID() gin.LogFormatter {
	return func(param gin.LogFormatterParams) string {
		// Initialize color strings for use in the terminal.
		var statusColor, methodColor, resetColor string
		if param.IsOutputColor() {
			// Set terminal color codes based on the log parameters.
			statusColor = param.StatusCodeColor()
			methodColor = param.MethodColor()
			resetColor = param.ResetColor()
		}

		// Truncate latency to the nearest second if it's greater than a minute.
		if param.Latency > time.Minute {
			param.Latency = param.Latency - (param.Latency % time.Second)
		}

		// Format the log message with the request information.
		return fmt.Sprintf("%s%3d%s - [%s] \"%v %s%s%s %s\" %s",
			statusColor, param.StatusCode, resetColor,
			param.ClientIP,
			param.Latency,
			methodColor, param.Method, resetColor,
			param.Path,
			param.ErrorMessage,
		)
	}
}
