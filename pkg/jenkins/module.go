package jenkins

import (
	"github.com/mgufrone/jenkins-slackbot/pkg/env"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

var Module = fx.Module("jenkins",
	fx.Decorate(func(logger logrus.FieldLogger) logrus.FieldLogger {
		return logger.WithField("component", "jenkins_http")
	}),
	fx.Provide(NewClient),
	fx.Invoke(func() error {
		return env.Requires(
			"JENKINS_URL",
			"JENKINS_USERNAME",
			"JENKINS_PASSWORD",
		)
	}),
)
