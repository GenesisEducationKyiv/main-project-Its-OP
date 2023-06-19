package application

import (
	"btcRate/domain"
	"fmt"
)

type CoinService struct {
	coinClient      domain.ICoinClient
	emailClient     domain.IEmailClient
	emailRepository domain.IEmailRepository
}

func NewCoinService(client domain.ICoinClient, emailClient domain.IEmailClient, emailRepository domain.IEmailRepository) *CoinService {
	return &CoinService{coinClient: client, emailClient: emailClient, emailRepository: emailRepository}
}

func (coinService *CoinService) GetCurrentRate(currency string, coin string) (*domain.Price, error) {
	rate, time, err := coinService.coinClient.GetRate(currency, coin)

	if err != nil {
		return nil, err
	}

	return &domain.Price{
		Amount:    rate,
		Currency:  currency,
		Timestamp: time,
	}, nil
}

func (coinService *CoinService) Subscribe(email string) error {
	err := coinService.emailRepository.AddEmail(email)
	if err != nil {
		return err
	}

	err = coinService.emailRepository.Save()
	if err != nil {
		return err
	}

	return nil
}

func (coinService *CoinService) SendRateEmails(currency string, coin string) {
	emails := coinService.emailRepository.GetAll()

	currentPrice, err := coinService.GetCurrentRate(currency, coin)
	if err != nil {
		return
	}

	htmlTemplate := `<p><strong>Amount:</strong> %f</p>
	<p><strong>Currency:</strong> %s<p>
	<p><strong>Timestamp:</strong> %s<p>`
	htmlBody := fmt.Sprintf(htmlTemplate, currentPrice.Amount, currentPrice.Currency, currentPrice.Timestamp.Format("02-01-06 15:04:05.999 Z0700"))

	coinService.emailClient.Send(emails, htmlBody)
}
