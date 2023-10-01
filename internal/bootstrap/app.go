package bootstrap

import (
	"github.com/mgufrone/jenkins-slackbot/pkg/logger"
	"go.uber.org/fx"
)

// bootstrap all core components required to run a certain service

var Core = fx.Module("core",
	logger.Module,
)
