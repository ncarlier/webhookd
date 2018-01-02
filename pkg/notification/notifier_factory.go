package notification

import (
	"errors"
	"os"
)

// Notifier is able to send a notification.
type Notifier interface {
	Notify(subject string, text string, attachfile string)
}

// NotifierFactory creates a notifier regarding the configuration.
func NotifierFactory() (Notifier, error) {
	notifier := os.Getenv("APP_NOTIFIER")
	switch notifier {
	case "http":
		return newHTTPNotifier(), nil
	case "smtp":
		return newSMTPNotifier(), nil
	default:
		if notifier == "" {
			return nil, errors.New("notification provider not configured")
		}
		return nil, errors.New("unknown notification provider: " + notifier)
	}
}
