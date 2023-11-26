package options

import (
	"github.com/skeleton1231/gotal/internal/pkg/server"
)

type RateLimitOptions struct {
	RequestsPerSecond float64                `json:"requests-per-second" mapstructure:"requests-per-second"`
	BurstSize         int                    `json:"burst-size" mapstructure:"burst-size"`
	CustomLimits      map[string]interface{} `json:"custom-limits" mapstructure:"custom-limits"`
}

func NewRateLimitOptions() *RateLimitOptions {
	defaults := server.NewConfig()
	return &RateLimitOptions{
		RequestsPerSecond: defaults.RateLimit.RequsetPerSecond,
		BurstSize:         defaults.RateLimit.Burst,
	}
}

func (r *RateLimitOptions) ApplyTo(c *server.Config) error {
	// 将默认速率限制设置应用到 server.Config
	// c.RateLimit.RequsetPerSecond = r.RequestsPerSecond
	// c.RateLimit.Burst = r.BurstSize

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
