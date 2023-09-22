package slack

import (
	"github.com/mgufrone/jenkins-slackbot/pkg/env"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

var Module = fx.Module("slack_socket",
	fx.Decorate(func(lg logrus.FieldLogger) logrus.FieldLogger {
		return lg.WithField("component", "slack_manager")
	}),
	fx.Provide(
		slackApi,
		socketClient,
		NewSocketManager,
	),
	fx.Invoke(
		func() error {
			return env.Requires("SLACK_APP_TOKEN", "SLACK_BOT_TOKEN")
		},
		func(manager *SocketManager) {},
	),
)
