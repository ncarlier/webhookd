package notification

import (
	"log/slog"
)

// Notifier is able to send a notification.
type Notifier interface {
	Notify(result HookResult) error
}

var notifier Notifier

// Notify is the global method to notify hook result
func Notify(result HookResult) {
	if notifier == nil {
		return
	}
	if err := notifier.Notify(result); err != nil {
		slog.Error("unable to send notification", "webhook", result.Name(), "id", result.ID(), "err", err)
	}
}

// Init creates the notifier singleton regarding the URI.
func Init(uri string) (err error) {
	notifier, err = NewNotifier(uri)
	return err
}
