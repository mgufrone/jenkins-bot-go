package routes

import (
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"

	"mgufrone.dev/job-tracking/app/http/controllers"
)

func Api() {
	userController := controllers.NewUserController()
	facades.Route().Prefix("/users").Group(func(router route.Router) {
		router.Get("/{id}", userController.Show)
	})
}
