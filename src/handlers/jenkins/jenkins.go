package jenkins

import (
	"context"
	"encoding/json"
	"github.com/mgufrone/jenkins-slackbot/pkg/jenkins"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"strconv"
	"strings"
)

type Jenkins struct {
	cli    *jenkins.Client
	logger logrus.FieldLogger
}

func NewJenkins(cli *jenkins.Client, logger logrus.FieldLogger) *Jenkins {
	return &Jenkins{cli: cli, logger: logger}
}

func (j *Jenkins) When(event socketmode.Event) bool {
	if event.Type != socketmode.EventTypeInteractive {
		return false
	}
	var pl slack.InteractionCallback
	by, _ := event.Request.Payload.MarshalJSON()
	if err := json.Unmarshal(by, &pl); err != nil {
		return false
	}
	switch pl.Type {
	case slack.InteractionTypeInteractionMessage:
		return true
	default:
		return false
	}
}

func (j *Jenkins) Run(api *slack.Client, ws *socketmode.Client, evt socketmode.Event) error {
	var pl slack.InteractionCallback
	by, _ := evt.Request.Payload.MarshalJSON()
	if err := json.Unmarshal(by, &pl); err != nil {
		return err
	}
	splits := strings.Split(pl.CallbackID, ":")
	action := pl.ActionCallback.AttachmentActions[0].Value
	jobName, buildIDStr := splits[0], splits[1]
	if strings.Contains(jobName, "/") {
		jobName = strings.ReplaceAll(jobName, "/", "/job/")
	}
	buildID, _ := strconv.Atoi(buildIDStr)
	logger := j.logger.WithField("job", jobName).WithField("buildID", buildID)
	logger.Infof("acknowledging action")
	ws.Ack(*evt.Request, map[string]interface{}{
		"text": ":gh-loading: Submitting your action",
	})
	submitter := pl.User.ID
	logger.Infof("inquire pending action info")
	act, err := j.cli.GetPendingAction(jobName, buildID)
	if err != nil {
		return err
	}
	defer func() {
		logger.Infof("action %s submitted", action)
	}()
	logger.Infof("submitting action: %s", action)
	if action == "approve" {
		return j.cli.Approve(context.TODO(), act, submitter)
	}
	return j.cli.Abort(context.TODO(), act)
}

func (j *Jenkins) Name() string {
	return "jenkins_slack_connector"
}
