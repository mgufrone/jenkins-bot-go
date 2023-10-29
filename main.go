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
		if err := facades.Route().Run("0.0.0.0:3000"); err != nil {
			facades.Log().Errorf("Route run error: %v", err)
		}
	}()

	select {}
}
