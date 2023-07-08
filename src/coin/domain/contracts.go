package domain

type ICoinService interface {
	GetCurrentRate(currency string, coin string) (*Price, error)
	SendRateEmails(currency string, coin string) error
}
