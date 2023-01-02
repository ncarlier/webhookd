package all

import (
	// activate HTTP notifier
	_ "github.com/ncarlier/webhookd/pkg/notification/http"
	// activate SMTP notifier
	_ "github.com/ncarlier/webhookd/pkg/notification/smtp"
)
