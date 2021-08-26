package main

import (
	"context"
	"errors"
	"github.com/mgufrone/jenkins-slackbot/handlers"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"go.uber.org/fx"
	"log"
	"os"
)

type InHandlers struct {
	fx.In
	Runner []*handlers.EventRegistrar `group:"interaction"`
}

func slackApi() *slack.Client {
	return slack.New(
		os.Getenv("SLACK_BOT_TOKEN"),
		slack.OptionAppLevelToken(os.Getenv("SLACK_APP_TOKEN")),
		slack.OptionDebug(os.Getenv("APP_ENV") != "production"),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
	)
}
func socketClient(apiClient *slack.Client) *socketmode.Client {
	return socketmode.New(
		apiClient,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "sm: ", log.Lshortfile|log.LstdFlags)),
	)
}

func testAuth(apiClient *slack.Client) (*slack.AuthTestResponse, error) {
	return apiClient.AuthTest()
}

func start(sockCli *socketmode.Client, webCli *slack.Client, lifecycle fx.Lifecycle, out InHandlers) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				for event := range sockCli.Events {
					for _, run := range out.Runner {
						if event.Type == run.EventType {
							eventPayload, _ := event.Data.(slack.InteractionCallback)
							if eventPayload.Type == run.InteractionType {
								if err := run.Handler(webCli, sockCli, event.Request); err != nil {
									log.Print("error", err)
								}
							}
						}
					}
				}
			}()
			return sockCli.RunContext(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}
func checkRequirements() error {
	if os.Getenv("JENKINS_USERNAME") == "" {
		return errors.New("require jenkins username")
	}
	if os.Getenv("JENKINS_PASSWORD") == "" {
		return errors.New("require jenkins password")
	}
	if os.Getenv("SLACK_BOT_TOKEN") == "" {
		return errors.New("require slack bot token")
	}
	if os.Getenv("SLACK_APP_TOKEN") == "" {
		return errors.New("require slack bot token")
	}
	return nil
}
func app() *fx.App {
	return fx.New(
		fx.Provide(
			slackApi,
			socketClient,
			testAuth,
			fx.Annotated{
				Group:  "interaction",
				Target: handlers.NewJenkinsHandler,
			},
		),
		fx.Invoke(
			checkRequirements,
			start,
		),
	)
}

func main() {
	a := app()
	if err := a.Start(context.Background()); err != nil {
		panic(err)
	}
}
