package globals

import (
	"context"
	"errors"
)

func SendNotification(ctx context.Context, notif Notification, dests []NotificationDestination) error {
	var errs []error
	for _, dest := range dests {
		err := dest.Send(ctx, notif)
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}
