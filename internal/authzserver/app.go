package authzserver

import (
	"github.com/skeleton1231/gotal/internal/authzserver/config"
	"github.com/skeleton1231/gotal/internal/authzserver/options"
	"github.com/skeleton1231/gotal/pkg/app"
	"github.com/skeleton1231/gotal/pkg/log"
)

const commandDesc = `Authorization server information`

func NewApp(basename string) *app.App {
	opts := options.NewOptions()
	application := app.NewApp("GoTal Authorization Server",
		basename,
		app.WithOptions(opts),
		app.WithDescription(commandDesc),
		app.WithDefaultValidArgs(),
		app.WithRunFunc(run(opts)),
	)
	return application
}

func run(opts *options.Options) app.RunFunc {
	return func(basename string) error {
		log.Init(opts.Log)
		defer log.Flush()

		cfg, err := config.CreateConfigFromOptions(opts)
		if err != nil {
			return err
		}

		// Run server
		return Run(cfg)
	}
}
