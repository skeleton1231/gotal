package options

import (
	"encoding/json"

	genericOptions "github.com/skeleton1231/gotal/internal/pkg/options"
	"github.com/skeleton1231/gotal/internal/pkg/server"
	"github.com/skeleton1231/gotal/pkg/log"
	"github.com/skeleton1231/gotal/pkg/util/flag"
)

// Options struct defines the configuration options for the server.
type Options struct {
	// RPCServer is the address of the RPC server.
	RPCServer string `json:"rpcserver" mapstructure:"rpcserver"`

	// ClientCA represents the file path of the client certificate authority.
	ClientCA string `json:"client-ca-file" mapstructure:"client-cat-file"`

	// GenericServerOptions holds the options for running a generic server.
	GenericServerOptions *genericOptions.ServerRunOptions `json:"server" mapstructure:"server"`

	// InsecureServing options for running the server without TLS.
	InsecureServing *genericOptions.InsecureServingOptions `json:"insecure" mapstructure:"insecure"`

	// SecureServing options for running the server with TLS.
	SecureServing *genericOptions.SecureServingOptions `json:"secure" mapstructure:"secure"`

	// RedisOptions specifies the options for connecting to Redis.
	RedisOptions *genericOptions.RedisOptions `json:"redis" mapstructure:"redis"`

	// FeatureOptions holds configurations for specific features.
	FeatureOptions *genericOptions.FeatureOptions `json:"feature" mapstructure:"featrue"`

	// Log holds the logging configuration.
	Log *log.Options `json:"log" mapstucture:"log"`
}

// NewOptions initializes a new Options object with default values.
func NewOptions() *Options {
	o := Options{
		RPCServer:            "127.0.0.1:8081",
		ClientCA:             "",
		GenericServerOptions: genericOptions.NewServerRunOptions(),
		InsecureServing:      genericOptions.NewInsecureServingOptions(),
		SecureServing:        genericOptions.NewSecureServingOptions(),
		RedisOptions:         genericOptions.NewRedisOptions(),
		FeatureOptions:       genericOptions.NewFeatureOptions(),
		Log:                  log.NewOptions(),
	}

	return &o
}

// ApplyTo applies the options to the given server configuration.
func (o *Options) ApplyTo(c *server.Config) error {
	// Implementation for applying options to config can be added here.
	return nil
}

// Flags defines and returns a set of flags for configuring the server.
func (o *Options) Flags() (fss flag.NamedFlagSets) {
	// Add flags for each option category.
	o.GenericServerOptions.AddFlags(fss.FlagSet("generic"))
	o.InsecureServing.AddFlags(fss.FlagSet("Insecure serving"))
	o.SecureServing.AddFlags(fss.FlagSet("Secure serving"))
	o.FeatureOptions.AddFlags(fss.FlagSet("features"))
	o.Log.AddFlags(fss.FlagSet("logs"))
	o.RedisOptions.AddFlags(fss.FlagSet("redis"))

	// Add miscellaneous flags.
	fs := fss.FlagSet("misc")
	fs.StringVar(&o.RPCServer, "rpcserver", o.RPCServer, "authorization rpc server")
	fs.StringVar(&o.ClientCA, "client-ca-file", o.ClientCA, "client certificate")
	return fss
}

// String returns a JSON representation of the options.
func (o *Options) String() string {
	data, _ := json.Marshal(o)
	return string(data)
}

// Complete finalizes the configuration of the secure serving options.
func (o *Options) Complete() error {
	// Complete the configuration for secure serving.
	return o.SecureServing.Complete()
}
