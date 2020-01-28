package notification

import (
	"fmt"
	"log"
	"net/smtp"
	"net/url"
	"strings"

	"github.com/ncarlier/webhookd/pkg/logger"
	"github.com/ncarlier/webhookd/pkg/model"
)

// SMTPNotifier is able to send notifcation to a email destination.
type SMTPNotifier struct {
	Host         string
	From         string
	To           string
	PrefixFilter string
}

func newSMTPNotifier(uri *url.URL) *SMTPNotifier {
	logger.Info.Println("using SMTP notification system: ", uri.Opaque)
	return &SMTPNotifier{
		Host:         getValueOrAlt(uri.Query(), "smtp", "localhost:25"),
		From:         getValueOrAlt(uri.Query(), "from", "noreply@nunux.org"),
		To:           uri.Opaque,
		PrefixFilter: getValueOrAlt(uri.Query(), "prefix", "notify:"),
	}
}

// Notify send a notification to a email destination.
func (n *SMTPNotifier) Notify(work *model.WorkRequest) error {
	// Get email body
	payload := work.GetLogContent(n.PrefixFilter)
	if strings.TrimSpace(payload) == "" {
		// Nothing to notify, abort
		return nil
	}

	// Buidl subject
	var subject string
	if work.Status == model.Success {
		subject = fmt.Sprintf("Webhook %s#%d SUCCESS.", work.Name, work.ID)
	} else {
		subject = fmt.Sprintf("Webhook %s#%d FAILED.", work.Name, work.ID)
	}

	// Connect to the remote SMTP server.
	c, err := smtp.Dial(n.Host)
	if err != nil {
		return err
	}

	// Set the sender and recipient first
	if err := c.Mail(n.From); err != nil {
		return err
	}
	if err := c.Rcpt(n.To); err != nil {
		log.Println(err)
		return err
	}

	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(wc, "Subject: %s\r\n\r\n%s\r\n\r\n", subject, payload)
	if err != nil {
		return err
	}
	err = wc.Close()
	if err != nil {
		return err
	}

	logger.Info.Printf("job %s#%d notification sent to %s\n", work.Name, work.ID, n.To)

	// Send the QUIT command and close the connection.
	return c.Quit()
}
