package slack_handler

import (
	"github.com/mgufrone/jenkins-slackbot/internal/services/slack"
	"go.uber.org/fx"
)

var Module = fx.Module("slack_handlers",
	fx.Provide(
		slack.AsSocketSubscriber(newJenkinsHandler),
	),
)
