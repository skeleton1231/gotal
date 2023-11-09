// Copyright 2023 Talhuang<talhuang1231@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package options contains flags and options for initializing an apiserver
package options

import (
	"encoding/json"

	"github.com/skeleton1231/gotal/pkg/util/flag"

	"github.com/skeleton1231/gotal/internal/pkg/options"
	"github.com/skeleton1231/gotal/internal/pkg/server"
	"github.com/skeleton1231/gotal/pkg/util/common"
)

type Options struct {
	GenericServerRunOptions *options.ServerRunOptions       `json:"server"   mapstructure:"server"`
	GRPCOptions             *options.GRPCOptions            `json:"grpc"     mapstructure:"grpc"`
	InsecureServing         *options.InsecureServingOptions `json:"insecure" mapstructure:"insecure"`
	SecureServing           *options.SecureServingOptions   `json:"secure"   mapstructure:"secure"`
	MySQLOptions            *options.MySQLOptions           `json:"mysql"    mapstructure:"mysql"`
	RedisOptions            *options.RedisOptions           `json:"redis"    mapstructure:"redis"`
	JwtOptions              *options.JwtOptions             `json:"jwt"      mapstructure:"jwt"`
	FeatureOptions          *options.FeatureOptions         `json:"feature"  mapstructure:"feature"`
	RateLimitOptions        *options.RateLimitOptions       `json:"ratelimit"  mapstructure:"ratelimit"`
}

func NewOptions() *Options {

	return &Options{
		GenericServerRunOptions: options.NewServerRunOptions(),
		GRPCOptions:             options.NewGRPCOptions(),
		InsecureServing:         options.NewInsecureServingOptions(),
		SecureServing:           options.NewSecureServingOptions(),
		MySQLOptions:            options.NewMySQLOptions(),
		RedisOptions:            options.NewRedisOptions(),
		JwtOptions:              options.NewJwtOptions(),
		FeatureOptions:          options.NewFeatureOptions(),
		RateLimitOptions:        options.NewRateLimitOptions(),
	}
}

func (o *Options) ApplyTo(c *server.Config) error {
	return nil
}

func (o *Options) Flags() (fss flag.NamedFlagSets) {
	o.GenericServerRunOptions.AddFlags(fss.FlagSet("generic"))
	o.JwtOptions.AddFlags(fss.FlagSet("jwt"))
	o.GRPCOptions.AddFlags(fss.FlagSet("grpc"))
	o.MySQLOptions.AddFlags(fss.FlagSet("mysql"))
	o.RedisOptions.AddFlags(fss.FlagSet("redis"))
	o.FeatureOptions.AddFlags(fss.FlagSet("features"))
	o.InsecureServing.AddFlags(fss.FlagSet("insecure serving"))
	o.SecureServing.AddFlags(fss.FlagSet("secure serving"))
	o.RateLimitOptions.AddFlags(fss.FlagSet("ratelimit"))
	return fss
}

func (o *Options) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}

// Complete set default Options.
func (o *Options) Complete() error {
	if o.JwtOptions.Key == "" {
		key, err := common.GenerateSecretKey(32) // Adjust the length as needed
		if err != nil {
			return err
		}
		o.JwtOptions.Key = key
	}

	return o.SecureServing.Complete()
}
