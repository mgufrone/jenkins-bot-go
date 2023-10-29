package providers

import (
	"github.com/goravel/framework/contracts/event"
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/facades"
	"mgufrone.dev/job-tracking/packages/jenkins/listeners"
	"mgufrone.dev/job-tracking/packages/slack/events"
	listeners2 "mgufrone.dev/job-tracking/packages/slack/listeners"
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
