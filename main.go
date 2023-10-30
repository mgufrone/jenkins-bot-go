package main

import (
	"github.com/goravel/framework/facades"

	"mgufrone.dev/job-tracking/bootstrap"
)

func main() {
	// This bootstraps the framework and gets it ready for use.
	bootstrap.Boot()

	//Start http server by facades.Route().
	go func() {
		if facades.Config().GetString("app.mode") == "aio" {
			facades.Artisan().Run([]string{".", "artisan", "slack:socket"}, false)
		}
	}()
	go func() {
		if err := facades.Route().Run(); err != nil {
			facades.Log().Errorf("Route run error: %v", err)
		}
	}()

	select {}
}
