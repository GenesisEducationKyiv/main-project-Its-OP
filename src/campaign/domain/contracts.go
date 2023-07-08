package domain

type ICampaignService interface {
	Subscribe(email string) error
	SendRateEmails() error
}
