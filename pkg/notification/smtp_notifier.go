package notification

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"net/url"
	"strings"
	"time"

	"github.com/ncarlier/webhookd/pkg/logger"
	"github.com/ncarlier/webhookd/pkg/model"
)

// SMTPNotifier is able to send notification to a email destination.
type SMTPNotifier struct {
	Host         string
	Username     string
	Password     string
	Conn         string
	From         string
	To           string
	PrefixFilter string
}

func newSMTPNotifier(uri *url.URL) *SMTPNotifier {
	logger.Info.Println("using SMTP notification system:", uri.Opaque)
	q := uri.Query()
	return &SMTPNotifier{
		Host:         getValueOrAlt(q, "smtp", "localhost:25"),
		Username:     getValueOrAlt(q, "username", ""),
		Password:     getValueOrAlt(q, "password", ""),
		Conn:         getValueOrAlt(q, "conn", "plain"),
		From:         getValueOrAlt(q, "from", "noreply@nunux.org"),
		To:           uri.Opaque,
		PrefixFilter: getValueOrAlt(q, "prefix", "notify:"),
	}
}

func (n *SMTPNotifier) buildEmailPayload(work *model.WorkRequest) string {
	// Get email body
	body := work.GetLogContent(n.PrefixFilter)
	if strings.TrimSpace(body) == "" {
		return ""
	}
	// Get email subject
	var subject string
	if work.Status == model.Success {
		subject = fmt.Sprintf("Webhook %s#%d SUCCESS.", work.Name, work.ID)
	} else {
		subject = fmt.Sprintf("Webhook %s#%d FAILED.", work.Name, work.ID)
	}
	// Build email headers
	headers := make(map[string]string)
	headers["From"] = n.From
	headers["To"] = n.To
	headers["Subject"] = subject

	// Build email payload
	payload := ""
	for k, v := range headers {
		payload += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	payload += "\r\n" + body
	return payload
}

// Notify send a notification to a email destination.
func (n *SMTPNotifier) Notify(work *model.WorkRequest) error {
	hostname, _, _ := net.SplitHostPort(n.Host)
	payload := n.buildEmailPayload(work)
	if payload == "" {
		// Nothing to notify, abort
		return nil
	}

	// Dial connection
	conn, err := net.DialTimeout("tcp", n.Host, 5*time.Second)
	if err != nil {
		return err
	}
	// Connect to SMTP server
	client, err := smtp.NewClient(conn, hostname)
	if err != nil {
		return err
	}

	if n.Conn == "tls" || n.Conn == "tls-insecure" {
		// TLS config
		tlsConfig := &tls.Config{
			InsecureSkipVerify: n.Conn == "tls-insecure",
			ServerName:         hostname,
		}
		if err := client.StartTLS(tlsConfig); err != nil {
			return err
		}
	}

	// Set auth if needed
	if n.Username != "" {
		if err := client.Auth(smtp.PlainAuth("", n.Username, n.Password, hostname)); err != nil {
			return err
		}
	}

	// Set the sender and recipient first
	if err := client.Mail(n.From); err != nil {
		return err
	}
	if err := client.Rcpt(n.To); err != nil {
		return err
	}

	// Send the email body.
	wc, err := client.Data()
	if err != nil {
		return err
	}

	_, err = wc.Write([]byte(payload))
	if err != nil {
		return err
	}
	err = wc.Close()
	if err != nil {
		return err
	}

	logger.Info.Printf("job %s#%d notification sent to %s\n", work.Name, work.ID, n.To)

	// Send the QUIT command and close the connection.
	return client.Quit()
}
