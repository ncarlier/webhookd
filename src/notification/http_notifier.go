package notification

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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

func (n HttpNotifier) Notify(subject string, text string, attachfile string) {
	log.Println("Sending notification '" + subject + "' to " + n.URL + " ...")
	data := make(url.Values)
	data.Set("from", n.From)
	data.Set("to", n.To)
	data.Set("subject", subject)
	data.Set("text", text)

	if attachfile != "" {
		file, err := os.Open(attachfile)
		if err != nil {
			log.Println(err)
			return
		}
		defer file.Close()

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("attachment", filepath.Base(attachfile))
		if err != nil {
			log.Println("Unable to create for file", err)
			return
		}
		_, err = io.Copy(part, file)

		for key, val := range data {
			_ = writer.WriteField(key, val[0])
		}

		err = writer.Close()
		if err != nil {
			log.Println("Unable to close writer", err)
			return
		}
		req, err := http.NewRequest("POST", n.URL, body)
		if err != nil {
			log.Println("Unable to post request", err)
		}
		defer req.Body.Close()
		req.Header.Set("Content-Type", writer.FormDataContentType())
		// Submit the request
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			log.Println("Unable to do the request", err)
			return
		}

		// Check the response
		if res.StatusCode != http.StatusOK {
			log.Println("bad status: %s", res.Status)
			return
		}
		log.Println("HTTP notification sended with attachment: ", attachfile)
	} else {
		resp, err := http.PostForm(n.URL, data)
		if err != nil {
			log.Println(err)
		}
		defer resp.Body.Close()
		log.Println("HTTP notification done")
	}
}
