package events

import "github.com/goravel/framework/contracts/event"

type InteractionSubmitted struct {
}

func (receiver *InteractionSubmitted) Handle(args []event.Arg) ([]event.Arg, error) {
	return args, nil
}
