// Copyright 2023 Talhuang<talhuang1231@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	maxAge = 24
)

func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		// Specify the actual domains in production
		AllowOrigins:  []string{"https://example.com", "https://api.example.com"},
		AllowMethods:  []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowHeaders:  []string{"Origin", "Authorization", "Content-Type", "Accept"},
		ExposeHeaders: []string{"Content-Length"},
		// Consider if you really need credentials with CORS
		AllowCredentials: false,
		AllowOriginFunc: func(origin string) bool {
			// Implement more complex logic if needed
			allowedOrigins := map[string]bool{
				"https://github.com": true,
				"https://google.com": true,
			}
			return allowedOrigins[origin]
		},
		// Adjust MaxAge according to how often your CORS policy may change
		MaxAge: maxAge * time.Hour,
	})
}
