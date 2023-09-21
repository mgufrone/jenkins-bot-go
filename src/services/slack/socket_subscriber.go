package slack

import (
	slack2 "github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"go.uber.org/fx"
)

type ISocketSubscriber interface {
	When(event socketmode.Event) bool
	Run(api *slack2.Client, ws *socketmode.Client, evt socketmode.Event) error
	Name() string
}

func AsSocketSubscriber(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(ISocketSubscriber)),
		fx.ResultTags(`group:"socket_subscribers"`),
	)
}
