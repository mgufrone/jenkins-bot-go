package slack

import (
	"github.com/mgufrone/jenkins-slackbot/pkg/env"
	"github.com/sirupsen/logrus"
	slack2 "github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"os"
)

type Logger struct {
	logger logrus.FieldLogger
}

func (l Logger) Output(i int, s string) error {
	if i == 1 {
		l.logger.Infoln(s)
		return nil
	}
	l.logger.Debugln(s)
	return nil
}

func slackApi(logger logrus.FieldLogger) *slack2.Client {
	slg := &Logger{logger}
	return slack2.New(
		env.Get("SLACK_BOT_TOKEN"),
		slack2.OptionAppLevelToken(os.Getenv("SLACK_APP_TOKEN")),
		slack2.OptionDebug(os.Getenv("APP_ENV") != "production"),
		slack2.OptionLog(slg),
	)
}
func socketClient(apiClient *slack2.Client, logger logrus.FieldLogger) *socketmode.Client {
	slg := &Logger{logger}
	return socketmode.New(
		apiClient,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(slg),
	)
}
