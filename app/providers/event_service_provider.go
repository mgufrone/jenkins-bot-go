package providers

import (
	"github.com/goravel/framework/contracts/event"
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/facades"
	"mgufrone.dev/jenkins-bot-go/packages/jenkins/listeners"
	"mgufrone.dev/jenkins-bot-go/packages/slack/events"
	listeners2 "mgufrone.dev/jenkins-bot-go/packages/slack/listeners"
)

type EventServiceProvider struct {
}

func (receiver *EventServiceProvider) Register(app foundation.Application) {
	facades.Event().Register(receiver.listen())
}

func (receiver *EventServiceProvider) Boot(app foundation.Application) {

}

func (receiver *EventServiceProvider) listen() map[event.Event][]event.Listener {
	return map[event.Event][]event.Listener{
		&events.InteractionSubmitted{}: {
			&listeners.OnInteractionSubmitted{},
		},
		&events.ResponseURL{}: {
			&listeners2.OnResponseURL{},
		},
	}
}
