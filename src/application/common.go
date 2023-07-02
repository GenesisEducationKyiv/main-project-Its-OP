package application

type MailBody struct {
	subject       string
	htmlContent   string
	receiverAlias string
}

type ICampaignService interface {
	Subscribe(email string) error
	SendEmails(body *MailBody) error
}
