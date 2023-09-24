package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mgufrone/jenkins-slackbot/pkg/env"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"net/http"
)

var Module = fx.Module("webserver",
	fx.Provide(
		func(lc fx.Lifecycle, logger logrus.FieldLogger) *gin.Engine {
			eng := gin.Default()
			srv := &http.Server{
				Addr:    fmt.Sprintf(":%d", env.Int("APP_PORT")),
				Handler: eng,
			}
			lc.Append(fx.StartStopHook(func(ctx context.Context) error {
				logger.Infoln("listening server at ", srv.Addr)
				return srv.ListenAndServe()
			}, func(ctx context.Context) error {
				return srv.Shutdown(ctx)
			}))
			return eng
		},
	), fx.Invoke(
		func(engine *gin.Engine, logger logrus.FieldLogger, route InRoute) {
			for _, v := range route.Routes {
				v.Mount(engine)
			}
		}, func(*gin.Engine) {},
	),
)
