package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/mgufrone/jenkins-slackbot/delivery/server/slack_handler"
	"github.com/mgufrone/jenkins-slackbot/internal/bootstrap"
	"github.com/mgufrone/jenkins-slackbot/internal/handlers/jenkins"
	slack2 "github.com/mgufrone/jenkins-slackbot/internal/services/slack"
	jenkins2 "github.com/mgufrone/jenkins-slackbot/pkg/jenkins"
	log "github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"os"
)

func main() {
	_ = godotenv.Load()
	a := fx.New(
		bootstrap.Core,
		slack2.Module,
		jenkins2.Module,
		fx.Provide(
			jenkins.NewJenkins,
		),
		slack_handler.Module,
	)
	if err := a.Start(context.TODO()); err != nil {
		log.Error(err)
		os.Exit(1)
	}
	<-a.Done()
}
