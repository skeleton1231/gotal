package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/time/rate"
)

// RateLimiterMiddleware returns a new rate limiter middleware.
func RateLimiter() gin.HandlerFunc {
	// Read configuration values using Viper
	rateLimit := viper.GetFloat64("ratelimit.requests-per-second") // Access nested configuration
	burstSize := viper.GetInt("ratelimit.burst-size")              // Access nested configuration

	// Convert rateLimit to a proper rate.Limit type
	r := rate.Limit(rateLimit)

	// Create a new rate limiter with the values retrieved from Viper
	limiter := rate.NewLimiter(r, burstSize)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message": "Too many requests",
			})
			return
		}
		c.Next()
	}
}
