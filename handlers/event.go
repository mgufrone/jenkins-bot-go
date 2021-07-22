package handlers

import (
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

type EventRegistrar struct {
	EventType       socketmode.EventType
	InteractionType slack.InteractionType
	Handler         Handler
}
type Handler func(web *slack.Client, socket *socketmode.Client, event *socketmode.Request) error
