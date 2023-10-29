package handlers

import (
	"fmt"
	"github.com/goravel/framework/contracts/http"
	"github.com/slack-go/slack"
	"mgufrone.dev/job-tracking/app/http/requests"
	facades2 "mgufrone.dev/job-tracking/packages/slack/facades"
)

type RequestApproval struct {
}

func NewRequestApproval() *RequestApproval {
	return &RequestApproval{}
}

func (r *RequestApproval) SendApproval(ctx http.Context) http.Response {
	var reqDeploy requests.RequestDeploy
	errors, err := ctx.Request().ValidateRequest(&reqDeploy)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]interface{}{
			"errors": err.Error(),
		})
	}
	if errors != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]interface{}{
			"errors": errors.All(),
		})
	}
	slackClient := facades2.SlackClient()
	attachment := slack.Attachment{
		Text:       reqDeploy.Text,
		Fallback:   reqDeploy.Text,
		CallbackID: fmt.Sprintf("%s:%d", reqDeploy.BuildName, reqDeploy.BuildNumber),
		Actions: []slack.AttachmentAction{
			{
				Name:  "prompt",
				Text:  "Abort",
				Type:  "button",
				Style: "danger",
				Value: "abort",
			},
			{
				Name:  "prompt",
				Text:  "Approve",
				Type:  "button",
				Style: "primary",
				Value: "approve",
			},
		},
	}
	target, ts, text, err := slackClient.SendMessageContext(ctx, reqDeploy.Channel, slack.MsgOptionAttachments(
		attachment,
	))
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}
	ctx.Response().Json(http.StatusOK, map[string]string{
		"channel": target,
		"ts":      ts,
		"text":    text,
	})
	return nil
}
