package authzserver

import "github.com/skeleton1231/gotal/internal/authzserver/config"

func Run(cfg *config.Config) error {
	server, err := createAuthzServer(cfg)
	if err != nil {
		return err
	}
	return server.PrepareRun().Run()
}
