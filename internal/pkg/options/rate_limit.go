package options

import (
	"github.com/skeleton1231/gotal/internal/pkg/server"
)

type RateLimitOptions struct {
	RequestsPerSecond float64              `json:"requests-per-second" mapstructure:"requests-per-second"`
	BurstSize         int                  `json:"burst-size" mapstructure:"burst-size"`
	CustomLimits      map[string]RateLimit `json:"custom-limits" mapstructure:"custom-limits"`
}

// RateLimit struct defines the settings for rate limiting.
type RateLimit struct {
	RequestsPerSecond float64 `json:"requests-per-second" mapstruct　ure:"requests-per-second"`
	BurstSize         int     `json:"burst-size" mapstructure:"burst-size"`
}

func NewRateLimitOptions() *RateLimitOptions {
	defaults := server.NewConfig()
	return &RateLimitOptions{
		RequestsPerSecond: defaults.RateLimit.RequsetPerSecond,
		BurstSize:         defaults.RateLimit.Burst,
		CustomLimits:      make(map[string]RateLimit),
	}
}

func (r *RateLimitOptions) ApplyTo(c *server.Config) error {

	return nil
}

// Validate checks and validates the user-provided parameters during program startup.
func (r *RateLimitOptions) Validate() []error {

	return nil
}

// AddFlags adds flags for a specific RateLimitOptions to the specified FlagSet.
// func (r *RateLimitOptions) AddFlags(fs *pflag.FlagSet) {
// 	fs.IntVar(&r.RequestsPerSecond, "ratelimit.requests-per-second", r.RequestsPerSecond, "Number of requests per second per user.")
// 	fs.IntVar(&r.BurstSize, "ratelimit.burst-size", r.BurstSize, "Maximum number of requests in a single burst.")
// }
