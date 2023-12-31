// Copyright 2023 Talhuang<talhuang1231@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package options

import (
	"fmt"
	"net"

	"github.com/spf13/pflag"
)

// RedisOptions defines options for redis cluster.
type RedisOptions struct {
	Host                  string   `json:"host"                     mapstructure:"host"                     description:"Redis service host address"`
	Port                  int      `json:"port"`
	Addrs                 []string `json:"addrs"                    mapstructure:"addrs"`
	Username              string   `json:"username"                 mapstructure:"username"`
	Password              string   `json:"password"                 mapstructure:"password"`
	Database              int      `json:"database"                 mapstructure:"database"`
	MasterName            string   `json:"master-name"              mapstructure:"master-name"`
	MaxIdle               int      `json:"optimisation-max-idle"    mapstructure:"optimisation-max-idle"`
	MaxActive             int      `json:"optimisation-max-active"  mapstructure:"optimisation-max-active"`
	Timeout               int      `json:"timeout"                  mapstructure:"timeout"`
	EnableCluster         bool     `json:"enable-cluster"           mapstructure:"enable-cluster"`
	UseSSL                bool     `json:"use-ssl"                  mapstructure:"use-ssl"`
	SSLInsecureSkipVerify bool     `json:"ssl-insecure-skip-verify" mapstructure:"ssl-insecure-skip-verify"`
}

// NewRedisOptions create a `zero` value instance.
func NewRedisOptions() *RedisOptions {
	return &RedisOptions{
		Host:                  "127.0.0.1",
		Port:                  6379,
		Addrs:                 []string{},
		Username:              "",
		Password:              "",
		Database:              0,
		MasterName:            "",
		MaxIdle:               2000,
		MaxActive:             4000,
		Timeout:               0,
		EnableCluster:         false,
		UseSSL:                false,
		SSLInsecureSkipVerify: false,
	}
}

// Validate verifies flags passed to RedisOptions.
func (o *RedisOptions) Validate() []error {
	errs := []error{}

	// Validate Host
	if o.Host == "" {
		errs = append(errs, fmt.Errorf("host cannot be empty"))
	}

	// Validate Port
	if o.Port < 1 || o.Port > 65535 {
		errs = append(errs, fmt.Errorf("port should be between 1 and 65535"))
	}

	// Validate Addrs
	for _, addr := range o.Addrs {
		host, port, err := net.SplitHostPort(addr)
		if err != nil || host == "" || port == "" {
			errs = append(errs, fmt.Errorf("invalid address format for %s, expected format: host:port", addr))
		}
	}

	// Validate Database when using Redis cluster
	if o.EnableCluster && o.Database != 0 {
		errs = append(errs, fmt.Errorf("database must be 0 when using Redis cluster"))
	}

	// Validate MaxIdle
	if o.MaxIdle <= 0 {
		errs = append(errs, fmt.Errorf("maxIdle should be positive"))
	}

	// Validate MaxActive
	if o.MaxActive < o.MaxIdle {
		errs = append(errs, fmt.Errorf("maxActive should be greater than or equal to maxIdle"))
	}

	// Validate Timeout
	if o.Timeout < 0 {
		errs = append(errs, fmt.Errorf("timeout should be non-negative"))
	}

	// Validate SSL flags if UseSSL is true
	if o.UseSSL && o.SSLInsecureSkipVerify {
		errs = append(errs, fmt.Errorf("using insecure SSL settings, make sure you understand the implications"))
	}

	return errs
}

// AddFlags adds flags related to redis storage for a specific APIServer to the specified FlagSet.
func (o *RedisOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Host, "redis.host", o.Host, "Hostname of your Redis server.")
	fs.IntVar(&o.Port, "redis.port", o.Port, "The port the Redis server is listening on.")
	fs.StringSliceVar(&o.Addrs, "redis.addrs", o.Addrs, "A set of redis address(format: 127.0.0.1:6379).")
	fs.StringVar(&o.Username, "redis.username", o.Username, "Username for access to redis service.")
	fs.StringVar(&o.Password, "redis.password", o.Password, "Optional auth password for Redis db.")

	fs.IntVar(&o.Database, "redis.database", o.Database, ""+
		"By default, the database is 0. Setting the database is not supported with redis cluster. "+
		"As such, if you have --redis.enable-cluster=true, then this value should be omitted or explicitly set to 0.")

	fs.StringVar(&o.MasterName, "redis.master-name", o.MasterName, "The name of master redis instance.")

	fs.IntVar(&o.MaxIdle, "redis.optimisation-max-idle", o.MaxIdle, ""+
		"This setting will configure how many connections are maintained in the pool when idle (no traffic). "+
		"Set the --redis.optimisation-max-active to something large, we usually leave it at around 2000 for "+
		"HA deployments.")

	fs.IntVar(&o.MaxActive, "redis.optimisation-max-active", o.MaxActive, ""+
		"In order to not over commit connections to the Redis server, we may limit the total "+
		"number of active connections to Redis. We recommend for production use to set this to around 4000.")

	fs.IntVar(&o.Timeout, "redis.timeout", o.Timeout, "Timeout (in seconds) when connecting to redis service.")

	fs.BoolVar(&o.EnableCluster, "redis.enable-cluster", o.EnableCluster, ""+
		"If you are using Redis cluster, enable it here to enable the slots mode.")

	fs.BoolVar(&o.UseSSL, "redis.use-ssl", o.UseSSL, ""+
		"If set, apiserver will assume the connection to Redis is encrypted. "+
		"(use with Redis providers that support in-transit encryption).")

	fs.BoolVar(&o.SSLInsecureSkipVerify, "redis.ssl-insecure-skip-verify", o.SSLInsecureSkipVerify, ""+
		"Allows usage of self-signed certificates when connecting to an encrypted Redis database.")
}
