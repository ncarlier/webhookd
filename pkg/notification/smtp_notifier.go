package notification

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

// SMTPNotifier is able to send notifcation to a email destination.
type SMTPNotifier struct {
	Host string
	From string
	To   string
}

func newSMTPNotifier() *SMTPNotifier {
	notifier := new(SMTPNotifier)
	notifier.Host = os.Getenv("APP_SMTP_NOTIFIER_HOST")
	if notifier.Host == "" {
		notifier.Host = "localhost:25"
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

// Notify send a notification to a email destination.
func (n *SMTPNotifier) Notify(subject string, text string, attachfile string) {
	log.Println("SMTP notification: ", subject)
	// Connect to the remote SMTP server.
	c, err := smtp.Dial(n.Host)
	if err != nil {
		log.Println(err)
		return
	}

	// Set the sender and recipient first
	if err := c.Mail(n.From); err != nil {
		log.Println(err)
		return
	}
	if err := c.Rcpt(n.To); err != nil {
		log.Println(err)
		return
	}

	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		log.Println(err)
		return
	}
	_, err = fmt.Fprintf(wc, text)
	if err != nil {
		log.Println(err)
		return
	}
	err = wc.Close()
	if err != nil {
		log.Println(err)
		return
	}

	// Send the QUIT command and close the connection.
	err = c.Quit()
	if err != nil {
		log.Fatal(err)
	}
}
