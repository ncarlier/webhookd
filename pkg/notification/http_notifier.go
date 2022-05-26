package notification

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/ncarlier/webhookd/pkg/logger"
)

type notifPayload struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Text  string `json:"text"`
	Error error  `json:"error,omitempty"`
}

// HTTPNotifier is able to send a notification to a HTTP endpoint.
type HTTPNotifier struct {
	URL          *url.URL
	PrefixFilter string
}

func newHTTPNotifier(uri *url.URL) *HTTPNotifier {
	logger.Info.Println("using HTTP notification system: ", uri.String())
	return &HTTPNotifier{
		URL:          uri,
		PrefixFilter: getValueOrAlt(uri.Query(), "prefix", "notify:"),
	}
}

// Notify send a notification to a HTTP endpoint.
func (n *HTTPNotifier) Notify(result HookResult) error {
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
	logger.Info.Printf("job %s#%d notification sent to %s\n", result.Name(), result.ID(), n.URL.String())
	return nil
}
