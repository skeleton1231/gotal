// Copyright 2023 Talhuang<talhuang1231@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package options

// Validate checks Options and return a slice of found errs.
func (o *Options) Validate() []error {
	var errs []error

	validators := []interface {
		Validate() []error
	}{
		o.GenericServerRunOptions,
		o.GRPCOptions,
		o.InsecureServing,
		o.SecureServing,
		o.MySQLOptions,
		o.RedisOptions,
		o.JwtOptions,
		o.FeatureOptions,
	}

	for _, validator := range validators {
		errs = append(errs, validator.Validate()...)
	}

	return errs
}
