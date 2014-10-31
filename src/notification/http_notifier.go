package notification

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type HttpNotifier struct {
	URL  string
	From string
	To   string
	User []string
}

func NewHttpNotifier() *HttpNotifier {
	notifier := new(HttpNotifier)
	notifier.URL = os.Getenv("APP_HTTP_NOTIFIER_URL")
	if notifier.URL == "" {
		log.Println("Unable to create HTTP notifier. APP_HTTP_NOTIFIER_URL not set.")
		return nil
	}
	user := os.Getenv("APP_HTTP_NOTIFIER_USER")
	if user != "" {
		notifier.User = strings.Split(user, ":")
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

func (n *HttpNotifier) Notify(subject string, text string, attachfile string) {
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

		mh := make(textproto.MIMEHeader)
		mh.Set("Content-Type", "application/x-gzip")
		mh.Set("Content-Disposition", fmt.Sprintf("form-data; name=\"attachment\"; filename=\"%s\"", filepath.Base(attachfile)))
		part, err := writer.CreatePart(mh)
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
			log.Println("Unable to close the gzip writer", err)
			return
		}
		req, err := http.NewRequest("POST", n.URL, body)
		if err != nil {
			log.Println("Unable to post request", err)
		}
		defer req.Body.Close()
		req.Header.Set("Content-Type", writer.FormDataContentType())

		if len(n.User) == 2 {
			req.SetBasicAuth(n.User[0], n.User[1])
		}

		// Submit the request
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			log.Println("Unable to do the request", err)
			return
		}

		// Check the response
		if res.StatusCode != http.StatusOK {
			log.Println("bad status: ", res.Status)
			log.Println(res.Body)
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
