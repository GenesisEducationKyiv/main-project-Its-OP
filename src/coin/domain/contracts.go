package domain

type ICoinService interface {
	GetCurrentRate(currency string, coin string) (*Price, error)
}
