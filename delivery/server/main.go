package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/mgufrone/jenkins-slackbot/delivery/server/slack_handler"
	"github.com/mgufrone/jenkins-slackbot/internal/bootstrap"
	"github.com/mgufrone/jenkins-slackbot/internal/handlers/jenkins"
	"github.com/mgufrone/jenkins-slackbot/pkg/slack"
	log "github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"os"
)

func main() {
	_ = godotenv.Load()
	a := fx.New(
		bootstrap.Core,
		fx.Provide(
			jenkins.NewJenkins,
			slack.AsSlackInteractionHandler(func(jenkins2 *jenkins.Jenkins) *jenkins.Jenkins {
				return jenkins2
			}),
		),
		//handlers.Module,
		slack_handler.Module,
		//server.Module,
	)
	if err := a.Start(context.TODO()); err != nil {
		log.Error(err)
		os.Exit(1)
	}
	<-a.Done()
}
