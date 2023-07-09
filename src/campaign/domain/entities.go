package domain

import "time"

type Rate struct {
	Amount    float64
	Currency  string
	Timestamp time.Time
}

type MailBody struct {
	ReceiverAlias string
	Subject       string
	HtmlContent   string
}
