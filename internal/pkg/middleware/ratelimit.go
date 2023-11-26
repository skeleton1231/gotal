package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/time/rate"
)

// RateLimit defines rate limit settings.
type RateLimit struct {
	RequestsPerSecond float64 `json:"requests-per-second" mapstruct　ure:"requests-per-second"`
	BurstSize         int     `json:"burst-size" mapstructure:"burst-size"`
}

func RateLimiter() gin.HandlerFunc {
	defaultRateLimit := viper.GetFloat64("ratelimit.requests-per-second")
	defaultBurstSize := viper.GetInt("ratelimit.burst-size")
	defaultLimiter := rate.NewLimiter(rate.Limit(defaultRateLimit), defaultBurstSize)

	customLimiters := make(map[string]*rate.Limiter)
	customLimitsConfig := viper.GetStringMap("ratelimit.custom-limits")
	for path, cfg := range customLimitsConfig {
		cfgMap := cfg.(map[string]interface{})

		// Determine the correct type for requests-per-second
		var reqPerSecond float64
		switch v := cfgMap["requests-per-second"].(type) {
		case int:
			reqPerSecond = float64(v)
		case float64:
			reqPerSecond = v
		default:
			logrus.Errorf("Invalid type for requests-per-second: %T", v)
			continue
		}

		limiter := rate.NewLimiter(rate.Limit(reqPerSecond), int(cfgMap["burst-size"].(int)))
		customLimiters[path] = limiter
		logrus.Infof("Custom limiter for path %s: requests-per-second: %v, burst-size: %v", path, reqPerSecond, cfgMap["burst-size"])
	}

	return func(c *gin.Context) {
		limiter := defaultLimiter
		if customLimiter, exists := customLimiters[c.Request.URL.Path]; exists {
			limiter = customLimiter
		}

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message": "Too many requests",
			})
			return
		}

		c.Next()
	}
}
