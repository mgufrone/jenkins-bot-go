package slack

import (
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/log"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"mgufrone.dev/jenkins-bot-go/packages/slack/commands"
)

const Binding = "slack.service"
const BindingClient = "slack.client"
const BindingSocket = "slack.socket"

var App foundation.Application

type ServiceProvider struct {
}

type slackLogger struct {
	logger log.Log
}

func (s slackLogger) Output(i int, s2 string) error {
	switch i {
	case 1:
		s.logger.Info(s2)
	default:
		s.logger.Debug(s2)
	}
	return nil
}

func (receiver *ServiceProvider) Register(app foundation.Application) {
	App = app

	app.Bind(Binding, func(app foundation.Application) (any, error) {
		return nil, nil
	})
	app.Bind(BindingClient, func(app foundation.Application) (any, error) {
		cfg := app.MakeConfig()
		slackBotToken := cfg.GetString("slack.bot_token")
		slackAppToken := cfg.GetString("slack.app_token")
		isDebug := cfg.GetBool("app.debug")
		logger := app.MakeLog()
		return slack.New(slackBotToken,
			slack.OptionAppLevelToken(slackAppToken),
			slack.OptionDebug(isDebug),
			slack.OptionLog(&slackLogger{logger}),
		), nil
	})
	app.Singleton(BindingSocket, func(app foundation.Application) (any, error) {
		cfg := app.MakeConfig()
		isDebug := cfg.GetBool("app.debug")
		logger := app.MakeLog()
		slackCli, _ := app.Make(BindingClient)
		return socketmode.New(slackCli.(*slack.Client), socketmode.OptionLog(&slackLogger{logger}), socketmode.OptionDebug(isDebug)), nil
	})
}

func (receiver *ServiceProvider) Boot(app foundation.Application) {
	cli, _ := app.Make(BindingSocket)
	eventManager := app.MakeEvent()
	app.Publishes("mgufrone.dev/jenkins-bot-go/packages/slack", map[string]string{
		"config/slack.go": app.ConfigPath("slack.go"),
	})
	app.Commands([]console.Command{
		commands.NewConsole(cli.(*socketmode.Client), eventManager),
	})
}
