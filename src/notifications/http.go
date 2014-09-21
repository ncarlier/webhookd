package notifications

import (
	"log"
	"net/http"
	"net/url"
	"os"
)

type HttpNotifier struct {
	URL  string
	From string
	To   string
}

func NewHttpNotifier() *HttpNotifier {
	notifier := new(HttpNotifier)
	notifier.URL = os.Getenv("APP_HTTP_NOTIFIER_URL")
	if notifier.URL == "" {
		log.Println("Unable to create HTTP notifier. APP_HTTP_NOTIFIER_URL not set.")
		return nil
	}
	notifier.From = os.Getenv("APP_NOTIFIER_FROM")
	if notifier.From == "" {
		notifier.From = "webhookd <noreply@nunux.org>"
	}
	notifier.To = os.Getenv("APP_NOTIFIER_TO")
	if notifier.To == "" {
		notifier.To = "hostmaster@nunux.org"
	}
	return notifier
}

func (n HttpNotifier) Notify(text string, subject string) {
	log.Println("HTTP notification: ", subject)
	data := make(url.Values)
	data.Set("from", n.From)
	data.Set("to", n.To)
	data.Set("subject", subject)
	data.Set("text", text)

	// Submit form
	resp, err := http.PostForm(n.URL, data)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
}
