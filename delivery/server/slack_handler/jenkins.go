package slack_handler

import (
	"context"
	"encoding/json"
	"github.com/mgufrone/jenkins-slackbot/internal/handlers/jenkins"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

type JenkinsHandler struct {
	handler *jenkins.Jenkins
	logger  logrus.FieldLogger
}

func newJenkinsHandler(handler *jenkins.Jenkins, logger logrus.FieldLogger) *JenkinsHandler {
	return &JenkinsHandler{handler: handler, logger: logger}
}

func (j *JenkinsHandler) When(event socketmode.Event) bool {
	if event.Type != socketmode.EventTypeInteractive {
		return false
	}
	var callback slack.InteractionCallback
	by, _ := event.Request.Payload.MarshalJSON()
	_ = json.Unmarshal(by, &callback)
	return j.handler.When(callback)
}

func (j *JenkinsHandler) Run(api *slack.Client, ws *socketmode.Client, evt socketmode.Event) error {
	var callback slack.InteractionCallback
	by, _ := evt.Request.Payload.MarshalJSON()
	_ = json.Unmarshal(by, &callback)
	msg, err := j.handler.Run(context.TODO(), callback)
	if err != nil {
		return err
	}
	ws.Ack(*evt.Request, msg)
	return nil
}

func (j *JenkinsHandler) Name() string {
	return "slack_jenkins_handler"
}
