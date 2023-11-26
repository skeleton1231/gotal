package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skeleton1231/gotal/pkg/log"
	"github.com/spf13/viper"
	"golang.org/x/time/rate"
)

// RateLimiter returns a middleware handler function for Gin.
func RateLimiter() gin.HandlerFunc {
	// Retrieve the default rate limit and burst size from the configuration.
	defaultRateLimit := viper.GetFloat64("ratelimit.requests-per-second")
	defaultBurstSize := viper.GetInt("ratelimit.burst-size")
	// Create a default limiter using these values.
	defaultLimiter := rate.NewLimiter(rate.Limit(defaultRateLimit), defaultBurstSize)

	// A map to hold custom limiters for specific paths.
	customLimiters := make(map[string]*rate.Limiter)
	// Retrieve custom limits configuration.
	customLimitsConfig := viper.GetStringMap("ratelimit.custom-limits")

	for path, cfg := range customLimitsConfig {
		cfgMap := cfg.(map[string]interface{})

		// Handle different potential types for requests-per-second.
		var reqPerSecond float64
		switch v := cfgMap["requests-per-second"].(type) {
		case int:
			reqPerSecond = float64(v)
		case float64:
			reqPerSecond = v
		default:
			// Log an error if the type is unexpected.
			log.Errorf("Invalid type for requests-per-second: %T", v)
			continue
		}

		// Create and store a custom limiter for the path.
		limiter := rate.NewLimiter(rate.Limit(reqPerSecond), int(cfgMap["burst-size"].(int)))
		customLimiters[path] = limiter
		// Log the custom limiter's details for debugging.
		log.Debugf("Custom limiter for path %s: requests-per-second: %v, burst-size: %v", path, reqPerSecond, cfgMap["burst-size"])
	}

	// Return the Gin middleware function.
	return func(c *gin.Context) {
		// Use the default limiter as a starting point.
		limiter := defaultLimiter
		// If there's a custom limiter for the current path, use it instead.
		if customLimiter, exists := customLimiters[c.Request.URL.Path]; exists {
			limiter = customLimiter
		}

		// Check if the request is allowed under the rate limit.
		if !limiter.Allow() {
			// If not, abort the request with a 429 Too Many Requests error.
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message": "Too many requests",
			})
			return
		}

		// If the request is allowed, continue to the next handler.
		c.Next()
	}
}
