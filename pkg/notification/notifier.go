package notification

import (
	"fmt"
	"net/url"
	"strings"

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

// Init creates a notifier regarding the URI.
func Init(uri string) error {
	if uri == "" {
		return nil
	}
	u, err := url.Parse(uri)
	if err != nil {
		return fmt.Errorf("invalid notification URL: %s", uri)
	}
	switch u.Scheme {
	case "mailto":
		notifier = newSMTPNotifier(u)
	case "http", "https":
		notifier = newHTTPNotifier(u)
	default:
		return fmt.Errorf("unable to create notifier: %v", err)
	}

	return nil
}

func getValueOrAlt(values url.Values, key, alt string) string {
	if val, ok := values[key]; ok {
		return strings.Join(val[:], ",")
	}
	return alt
}
