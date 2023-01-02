package notification

import (
	"github.com/ncarlier/webhookd/pkg/logger"
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
		logger.Error.Printf("unable to send notification for webhook %s#%d: %v\n", result.Name(), result.ID(), err)
	}
}

// Init creates the notifier singleton regarding the URI.
func Init(uri string) (err error) {
	notifier, err = NewNotifier(uri)
	return err
}
