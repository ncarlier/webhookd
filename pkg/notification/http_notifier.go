package notification

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/ncarlier/webhookd/pkg/logger"
)

// HTTPNotifier is able to send a notification to a HTTP endpoint.
type HTTPNotifier struct {
	URL  string
	From string
	To   string
	User []string
}

func newHTTPNotifier() *HTTPNotifier {
	notifier := new(HTTPNotifier)
	notifier.URL = os.Getenv("APP_HTTP_NOTIFIER_URL")
	if notifier.URL == "" {
		logger.Error.Println("Unable to create HTTP notifier. APP_HTTP_NOTIFIER_URL not set.")
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

// Notify send a notification to a HTTP endpoint.
func (n *HTTPNotifier) Notify(subject string, text string, attachfile string) {
	logger.Debug.Println("Sending notification '" + subject + "' to " + n.URL + " ...")
	data := make(url.Values)
	data.Set("from", n.From)
	data.Set("to", n.To)
	data.Set("subject", subject)
	data.Set("text", text)

	if attachfile != "" {
		file, err := os.Open(attachfile)
		if err != nil {
			logger.Error.Println("Unable to open notification attachment file", err)
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
			logger.Error.Println("Unable to create HTTP notification attachment", err)
			return
		}
		_, err = io.Copy(part, file)

		for key, val := range data {
			_ = writer.WriteField(key, val[0])
		}

		err = writer.Close()
		if err != nil {
			logger.Error.Println("Unable to close the gzip writer", err)
			return
		}
		req, err := http.NewRequest("POST", n.URL, body)
		if err != nil {
			logger.Error.Println("Unable to post HTTP notification", err)
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
			logger.Error.Println("Unable to do HTTP notification request", err)
			return
		}

		// Check the response
		if res.StatusCode != http.StatusOK {
			logger.Error.Println("HTTP notification bad response: ", res.Status)
			logger.Debug.Println(res.Body)
			return
		}
		logger.Info.Println("HTTP notification sent with attachment: ", attachfile)
	} else {
		req, err := http.NewRequest("POST", n.URL, bytes.NewBufferString(data.Encode()))
		if err != nil {
			logger.Error.Println("Unable to post HTTP notification request", err)
		}
		defer req.Body.Close()

		if len(n.User) == 2 {
			req.SetBasicAuth(n.User[0], n.User[1])
		}

		// Submit the request
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			logger.Error.Println("Unable to do the HTTP notification request", err)
			return
		}

		// Check the response
		if res.StatusCode != http.StatusOK {
			logger.Error.Println("HTTP notification bad response: ", res.Status)
			logger.Debug.Println(res.Body)
			return
		}
		logger.Info.Println("HTTP notification sent.")
	}
}
