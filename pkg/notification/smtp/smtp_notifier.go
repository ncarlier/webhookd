package smtp

import (
	"crypto/tls"
	"fmt"
	"log/slog"
	"net"
	"net/smtp"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ncarlier/webhookd/pkg/helper"
	"github.com/ncarlier/webhookd/pkg/notification"
)

// smtpNotifier is able to send notification to a email destination.
type smtpNotifier struct {
	Host         string
	Username     string
	Password     string
	Conn         string
	From         string
	To           string
	Subject      string
	PrefixFilter string
}

func newSMTPNotifier(uri *url.URL) (notification.Notifier, error) {
	slog.Info("using SMTP notification system", "uri", uri.Opaque)
	q := uri.Query()
	return &smtpNotifier{
		Host:         helper.GetValueOrAlt(q, "smtp", "localhost:25"),
		Username:     helper.GetValueOrAlt(q, "username", ""),
		Password:     helper.GetValueOrAlt(q, "password", ""),
		Conn:         helper.GetValueOrAlt(q, "conn", "plain"),
		From:         helper.GetValueOrAlt(q, "from", "noreply@nunux.org"),
		To:           uri.Opaque,
		Subject:      helper.GetValueOrAlt(uri.Query(), "subject", "[whd-notification] {name}#{id} {status}"),
		PrefixFilter: helper.GetValueOrAlt(uri.Query(), "prefix", "notify:"),
	}, nil
}

func (n *smtpNotifier) buildEmailPayload(result notification.HookResult) string {
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
func (n *smtpNotifier) Notify(result notification.HookResult) error {
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

	slog.Info("notification sent", "hook", result.Name(), "id", result.ID(), "to", n.To)

	// Send the QUIT command and close the connection.
	return client.Quit()
}

func buildSubject(template string, result notification.HookResult) string {
	subject := strings.ReplaceAll(template, "{name}", result.Name())
	subject = strings.ReplaceAll(subject, "{id}", strconv.FormatUint(uint64(result.ID()), 10))
	subject = strings.ReplaceAll(subject, "{status}", result.StatusLabel())
	return subject
}

func init() {
	notification.Register("mailto", newSMTPNotifier)
}
