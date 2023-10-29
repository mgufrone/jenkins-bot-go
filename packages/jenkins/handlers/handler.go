package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/goravel/framework/contracts/event"
	"github.com/goravel/framework/contracts/log"
	"github.com/goravel/framework/facades"
	"github.com/slack-go/slack"
	"mgufrone.dev/job-tracking/packages/jenkins/contracts"
	"mgufrone.dev/job-tracking/packages/slack/events"
	"strconv"
	"strings"
)

type Jenkins struct {
	svc    contracts.Jenkins
	logger log.Log
}

func NewJenkins(svc contracts.Jenkins, logger log.Log) *Jenkins {
	return &Jenkins{svc: svc, logger: logger}
}

func (j *Jenkins) Run(ctx context.Context, payload string) error {

	var callback slack.InteractionCallback

	if err := json.NewDecoder(strings.NewReader(payload)).Decode(&callback); err != nil {
		return err
	}
	if callback.Type != slack.InteractionTypeInteractionMessage {
		return nil
	}
	splits := strings.Split(callback.CallbackID, ":")
	action := callback.ActionCallback.AttachmentActions[0].Value
	jobName, buildIDStr := splits[0], splits[1]
	if strings.Contains(jobName, "/") {
		jobName = strings.ReplaceAll(jobName, "/", "/job/")
	}
	buildID, _ := strconv.Atoi(buildIDStr)
	logger := j.logger.In(fmt.Sprintf("jenkins.%s.%d", jobName, buildID))
	logger.Infof("acknowledging action")
	acknowledgeText := facades.Config().GetString("jenkins.acknowledgeText")
	by, _ := json.Marshal(callback)

	submitter := callback.User.ID
	logger.Debugf("inquire pending action info for %s:%d", jobName, buildID)
	act, err := j.svc.GetPendingAction(jobName, buildID)
	if err != nil {
		logger.Error(err)
		return err
	}
	logger.Debugf("submitting action: %s", action)
	defer func() {
		if err != nil {
			logger.Error(err)
			go facades.Event().Job(&events.ResponseURL{}, []event.Arg{
				{Type: "string", Value: string(by)},
				{Type: "string", Value: facades.Config().GetString("jenkins.text.error", "failed to submit action. Consult with the administrator. In the meantime, please take action manually")},
			}).Dispatch()
		}
	}()
	handlerMap := map[string]func() error{
		"approve": func() error {
			act.Submitter = submitter
			return j.svc.Approve(ctx, act)
		},
		"abort": func() error {
			return j.svc.Reject(ctx, act)
		},
	}
	go facades.Event().Job(&events.ResponseURL{}, []event.Arg{
		{Type: "string", Value: string(by)},
		{Type: "string", Value: acknowledgeText},
	}).Dispatch()
	if handler, ok := handlerMap[action]; ok {
		err = handler()
	}
	if err == nil {
		err = facades.Event().Job(&events.ResponseURL{}, []event.Arg{
			{Type: "string", Value: string(by)},
			{Type: "string", Value: fmt.Sprintf("%s by <@%s>", action, submitter)},
		}).Dispatch()
	}
	return err

}
