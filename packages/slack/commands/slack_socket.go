package commands

import (
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/contracts/event"
	"github.com/goravel/framework/facades"
	"github.com/slack-go/slack/socketmode"
	"mgufrone.dev/job-tracking/packages/slack/events"
)

type Console struct {
	socket     *socketmode.Client
	evtManager event.Instance
}

func NewConsole(socket *socketmode.Client, evtManager event.Instance) *Console {
	return &Console{socket: socket, evtManager: evtManager}
}

func (c *Console) Signature() string {
	return "slack:socket"
}

func (c *Console) Description() string {
	return "Run slack websocket listener"
}

func (c *Console) Extend() command.Extend {
	return command.Extend{}
}

func (c *Console) Handle(ctx console.Context) error {
	logger := facades.Log()
	go func() {
		for evt := range c.socket.Events {
			if evt.Type != socketmode.EventTypeInteractive {
				continue
			}
			by, _ := evt.Request.Payload.MarshalJSON()
			//_ = json.Unmarshal(by, &callback)
			logger.Info("publishing event", "interaction submitted")
			c.evtManager.Job(&events.InteractionSubmitted{}, []event.Arg{
				{Value: string(by), Type: "string"},
			}).Dispatch()
			c.socket.Ack(*evt.Request)
		}
	}()
	return c.socket.Run()
	//return
}
