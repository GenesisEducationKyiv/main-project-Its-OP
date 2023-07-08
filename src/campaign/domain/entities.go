package domain

type Rate struct {
	Amount   float64
	Currency string
}

type MailBody struct {
	Subject       string
	HtmlContent   string
	ReceiverAlias string
}
