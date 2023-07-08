package domain

import "time"

type Rate struct {
	Amount    float64
	Timestamp time.Time
	Currency  string
}

type MailBody struct {
	Subject       string
	HtmlContent   string
	ReceiverAlias string
}
