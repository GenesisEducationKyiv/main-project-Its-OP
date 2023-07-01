package domain

import "time"

type ICoinService interface {
	GetCurrentRate(currency string, coin string) (*Price, error)
	Subscribe(email string) error
	SendRateEmails(currency string, coin string) error
}

//go:generate mockery --name ICoinClient
type ICoinClient interface {
	GetRate(currency string, coin string) (float64, time.Time, error)
}

type IEmailRepository interface {
	AddEmail(email string) error
	GetAll() []string
	Save() error
}

type IEmailClient interface {
	Send(recipients []string, htmlContent string) error
}
