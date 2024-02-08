// Copyright 2023 Talhuang<talhuang1231@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package userservice

import "github.com/skeleton1231/gotal/internal/user_service/config"

// Run runs the specified APIServer. This should never exit.
func Run(cfg *config.Config) error {
	server, err := NewAPIServer(cfg)
	if err != nil {
		return err
	}

	return server.PrepareRun().Run()
}
