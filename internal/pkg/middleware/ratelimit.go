package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/time/rate"
)

// RateLimiterMiddleware returns a new rate limiter middleware.
func RateLimiter() gin.HandlerFunc {
	// Read configuration values using Viper
	rateLimit := viper.GetFloat64("rate_limit.requests-per-second") // Access nested configuration
	burstSize := viper.GetInt("rate_limit.burst-size")              // Access nested configuration

	logrus.Infof("rateLimit is %v", rateLimit)
	logrus.Infof("burstSize is %v", burstSize)

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
