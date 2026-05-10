package globals

import (
	"context"
)

type Bot interface {
	Run(context.Context) error
	Name() string
}

// discord webhook, text/email message
type NotificationDestination interface {
	Send(context.Context, Notification) error
}
type NotificationCache interface {
	Cache(Notification) error
	Reset() error
}
