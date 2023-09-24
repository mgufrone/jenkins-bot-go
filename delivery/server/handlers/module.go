package handlers

import (
	"github.com/mgufrone/jenkins-slackbot/pkg/server"
	"go.uber.org/fx"
)

var Module = fx.Module("http_routes",
	fx.Provide(
		server.AsRouteHandler(slackInteractionHandler, fx.ParamTags(`group:"slack_interaction_handlers"`)),
		server.AsRouteHandler(newSlackOauthHandler, fx.ParamTags()),
	),
)
