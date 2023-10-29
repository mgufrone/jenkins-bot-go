package events

import "github.com/goravel/framework/contracts/event"

type ResponseURL struct {
}

func (receiver *ResponseURL) Handle(args []event.Arg) ([]event.Arg, error) {
	return args, nil
}
