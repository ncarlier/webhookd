package notification

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ncarlier/webhookd/pkg/logger"
)

// SMTPNotifier is able to send notification to a email destination.
type SMTPNotifier struct {
	Host         string
	Username     string
	Password     string
	Conn         string
	From         string
	To           string
	Subject      string
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
		Subject:      getValueOrAlt(uri.Query(), "subject", "[whd-notification] {name}#{id} {status}"),
		PrefixFilter: getValueOrAlt(uri.Query(), "prefix", "notify:"),
	}
}

func (n *SMTPNotifier) buildEmailPayload(result HookResult) string {
	// Get email body
	body := result.Logs(n.PrefixFilter)
	if strings.TrimSpace(body) == "" {
		return ""
	}

	// Build email subject
	subject := buildSubject(n.Subject, result)

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
func (n *SMTPNotifier) Notify(result HookResult) error {
	hostname, _, _ := net.SplitHostPort(n.Host)
	payload := n.buildEmailPayload(result)
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

	logger.Info.Printf("job %s#%d notification sent to %s\n", result.Name(), result.ID(), n.To)

	// Send the QUIT command and close the connection.
	return client.Quit()
}

func buildSubject(template string, result HookResult) string {
	subject := strings.ReplaceAll(template, "{name}", result.Name())
	subject = strings.ReplaceAll(subject, "{id}", strconv.FormatUint(uint64(result.ID()), 10))
	subject = strings.ReplaceAll(subject, "{status}", result.StatusLabel())
	return subject
}
