package application

type MailBody struct {
	Subject       string
	HtmlContent   string
	ReceiverAlias string
}

type ICampaignService interface {
	Subscribe(email string) error
	SendEmails(body *MailBody) error
}
