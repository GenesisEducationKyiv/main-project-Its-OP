package domain

type ICampaignService interface {
	Subscribe(email string) error
	SendRateEmails(body *MailBody) error
}
