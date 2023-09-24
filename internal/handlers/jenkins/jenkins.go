package jenkins

import (
	"context"
	"github.com/mgufrone/jenkins-slackbot/pkg/jenkins"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"strconv"
	"strings"
)

type Jenkins struct {
	cli    *jenkins.Client
	logger logrus.FieldLogger
}

func (j *Jenkins) When(callback slack.InteractionCallback) bool {
	switch callback.Type {
	case slack.InteractionTypeInteractionMessage:
		return true
	default:
		return false
	}
}

func (j *Jenkins) Run(ctx context.Context, callback slack.InteractionCallback) (res slack.Message, err error) {
	splits := strings.Split(callback.CallbackID, ":")
	action := callback.ActionCallback.AttachmentActions[0].Value
	jobName, buildIDStr := splits[0], splits[1]
	if strings.Contains(jobName, "/") {
		jobName = strings.ReplaceAll(jobName, "/", "/job/")
	}
	buildID, _ := strconv.Atoi(buildIDStr)
	logger := j.logger.WithField("job", jobName).WithField("buildID", buildID)
	logger.Infof("acknowledging action")
	res = slack.Message{}
	res.Text = ":gh-loading: Submitting your action"
	res.ReplaceOriginal = true
	defer func() {
		submitter := callback.User.ID
		logger.Infof("inquire pending action info")
		act, err := j.cli.GetPendingAction(jobName, buildID)
		if err != nil {
			logger.Errorln(err)
			return
		}
		logger.Infof("submitting action: %s", action)
		if action == "approve" {
			err = j.cli.Approve(context.TODO(), act, submitter)
			if err != nil {
				logger.Errorln(err)
			}
			return
		}
		err = j.cli.Abort(context.TODO(), act)
		if err != nil {
			logger.Errorln(err)
		}
		logger.Infof("action %s submitted", action)
	}()
	return
}

func NewJenkins(cli *jenkins.Client, logger logrus.FieldLogger) *Jenkins {
	return &Jenkins{cli: cli, logger: logger}
}

func (j *Jenkins) Name() string {
	return "jenkins_slack_connector"
}
