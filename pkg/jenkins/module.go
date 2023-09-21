package jenkins

import (
	"github.com/mgufrone/jenkins-slackbot/pkg/env"
	"go.uber.org/fx"
)

var Module = fx.Module("jenkins",
	fx.Provide(NewClient),
	fx.Invoke(func() error {
		return env.Requires(
			"JENKINS_URL",
			"JENKINS_USERNAME",
			"JENKINS_PASSWORD",
		)
	}),
)
