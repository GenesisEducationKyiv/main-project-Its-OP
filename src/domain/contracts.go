package domain

type ICoinService interface {
	GetCurrentRate(currency string, coin string) (*Price, error)
	SendRateEmails(currency string, coin string) error
}

type ICampaignService interface {
	Subscribe(email string) error
	SendEmails(htmlBody string) error
}
