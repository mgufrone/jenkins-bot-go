package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/mgufrone/jenkins-slackbot/pkg/jenkins"
	"github.com/mgufrone/jenkins-slackbot/pkg/logger"
	jenkins2 "github.com/mgufrone/jenkins-slackbot/src/handlers/jenkins"
	slack2 "github.com/mgufrone/jenkins-slackbot/src/services/slack"
	"go.uber.org/fx"
	"log"
	"os"
)

func app() *fx.App {
	return fx.New(
		fx.NopLogger,
		logger.Module,
		jenkins.Module,
		fx.Provide(
			slack2.AsSocketSubscriber(jenkins2.NewJenkins),
		),
		slack2.Module,
	)
}

func main() {
	_ = godotenv.Load()
	a := app()

	if err := a.Start(context.Background()); err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	<-a.Done()
}
