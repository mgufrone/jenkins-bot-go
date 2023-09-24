package server

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type IHandler interface {
	Mount(srv gin.IRouter)
}

type InRoute struct {
	fx.In
	Routes []IHandler `group:"server_routes"`
}

func AsRouteHandler(f any, dependency fx.Annotation) any {
	return fx.Annotate(
		f,
		fx.As(new(IHandler)),
		fx.ResultTags(`group:"server_routes"`),
		dependency,
	)
}
