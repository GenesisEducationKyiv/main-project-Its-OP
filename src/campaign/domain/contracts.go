package domain

type ICampaignService interface {
	Subscribe(email string) error
	SendRateEmails(currency string, coin string) error
}
