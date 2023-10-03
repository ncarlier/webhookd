package http

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/ncarlier/webhookd/pkg/helper"
	"github.com/ncarlier/webhookd/pkg/notification"
)

type notifPayload struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Text  string `json:"text"`
	Error error  `json:"error,omitempty"`
}

// httpNotifier is able to send a notification to a HTTP endpoint.
type httpNotifier struct {
	URL          *url.URL
	PrefixFilter string
}

func newHTTPNotifier(uri *url.URL) (notification.Notifier, error) {
	slog.Info("using HTTP notification system ", "üri", uri.Opaque)
	return &httpNotifier{
		URL:          uri,
		PrefixFilter: helper.GetValueOrAlt(uri.Query(), "prefix", "notify:"),
	}, nil
}

// Notify send a notification to a HTTP endpoint.
func (n *httpNotifier) Notify(result notification.HookResult) error {
	payload := result.Logs(n.PrefixFilter)
	if strings.TrimSpace(payload) == "" {
		// Nothing to notify, abort
		return nil
	}

	notif := &notifPayload{
		ID:    strconv.FormatUint(result.ID(), 10),
		Name:  result.Name(),
		Text:  payload,
		Error: result.Err(),
	}
	notifJSON, err := json.Marshal(notif)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", n.URL.String(), bytes.NewBuffer(notifJSON))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	slog.Info("notification sent", "hook", result.Name(), "id", result.ID(), "to", n.URL.Opaque)
	return nil
}

func init() {
	notification.Register("http", newHTTPNotifier)
	notification.Register("https", newHTTPNotifier)
}
