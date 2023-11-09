package listeners

import (
	"encoding/json"
	"github.com/goravel/framework/contracts/event"
	slack2 "github.com/slack-go/slack"
	facades2 "mgufrone.dev/jenkins-bot-go/packages/slack/facades"
)

type OnResponseURL struct {
}

func (receiver *OnResponseURL) Signature() string {
	return "on_response_url"
}

func (receiver *OnResponseURL) Queue(args ...any) event.Queue {
	return event.Queue{
		Enable:     true,
		Connection: "",
		Queue:      "",
	}
}

func (receiver *OnResponseURL) Handle(args ...any) (err error) {
	svc := facades2.SlackClient()

	var msg slack2.InteractionCallback
	str := args[0].(string)
	if err = json.Unmarshal([]byte(str), &msg); err != nil {
		return
	}
	slackMsg := args[1].(string)
	_, _, _, err = svc.SendMessage(msg.Channel.ID,
		slack2.MsgOptionReplaceOriginal(msg.ResponseURL),
		slack2.MsgOptionAttachments(slack2.Attachment{Text: slackMsg, Fallback: slackMsg}),
	)
	return
}
