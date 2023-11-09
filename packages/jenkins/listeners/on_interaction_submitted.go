package listeners

import (
	"context"
	"github.com/goravel/framework/contracts/event"
	facades2 "mgufrone.dev/jenkins-bot-go/packages/jenkins/facades"
)

type OnInteractionSubmitted struct {
}

func (receiver *OnInteractionSubmitted) Signature() string {
	return "on_interaction_submitted"
}

func (receiver *OnInteractionSubmitted) Queue(args ...any) event.Queue {
	return event.Queue{
		Enable:     true,
		Connection: "",
		Queue:      "",
	}
}

func (receiver *OnInteractionSubmitted) Handle(args ...any) error {
	return facades2.Handler().Run(context.TODO(), args[0].(string))
}
