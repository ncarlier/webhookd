package notification

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/ncarlier/webhookd/pkg/logger"
	"github.com/ncarlier/webhookd/pkg/model"
)

// Notifier is able to send a notification.
type Notifier interface {
	Notify(work *model.WorkRequest) error
}

var notifier Notifier

// Notify is the global method to notify work
func Notify(work *model.WorkRequest) {
	if notifier == nil {
		return
	}
	if err := notifier.Notify(work); err != nil {
		logger.Error.Printf("unable to send notification for webhook %s#%d: %v\n", work.Name, work.ID, err)
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
