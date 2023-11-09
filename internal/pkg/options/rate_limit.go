package options

import (
	"github.com/skeleton1231/gotal/internal/pkg/server"
	"github.com/spf13/pflag"
)

// RateLimitOptions holds configuration for the rate limiting feature.
type RateLimitOptions struct {
	RequestsPerSecond int `json:"requests-per-second" mapstructure:"requests-per-second"`
	BurstSize         int `json:"burst-size" mapstructure:"burst-size"`
}

// NewRateLimitOptions creates a new RateLimitOptions with default values.
func NewRateLimitOptions() *RateLimitOptions {
	defaults := server.NewConfig()
	return &RateLimitOptions{
		RequestsPerSecond: int(defaults.RateLimit.TokensPerSecond), // Default requests per second
		BurstSize:         defaults.RateLimit.Burst,                // Default burst size
	}
}

// ApplyTo applies the current options to the provided server configuration.
// func (r *RateLimitOptions) ApplyTo(c *server.Config) error {
// 	c.RateLimit = &server.RateLimit{
// 		TokensPerSecond: r.RequestsPerSecond,
// 		Burst:           r.BurstSize,
// 	}
// 	return nil
// }

// Validate checks and validates the user-provided parameters during program startup.
func (r *RateLimitOptions) Validate() []error {

	return nil
}

// AddFlags adds flags for a specific RateLimitOptions to the specified FlagSet.
func (r *RateLimitOptions) AddFlags(fs *pflag.FlagSet) {
	fs.IntVar(&r.RequestsPerSecond, "ratelimit.requests-per-second", r.RequestsPerSecond, "Number of requests per second per user.")
	fs.IntVar(&r.BurstSize, "ratelimit.burst-size", r.BurstSize, "Maximum number of requests in a single burst.")
}
