package application

import (
	"btcRate/domain"
	"fmt"
	"golang.org/x/exp/slices"
	"time"
)

type ICoinClient interface {
	GetRate(currency string, coin string) (float64, time.Time, error)
}

type CoinService struct {
	supportedCurrencies []string
	supportedCoins      []string
	coinClient          ICoinClient
	emailClient         domain.IEmailClient
	emailRepository     domain.IEmailRepository
}

func NewCoinService(supportedCurrencies []string, supportedCoins []string, client ICoinClient, emailClient domain.IEmailClient, emailRepository domain.IEmailRepository) *CoinService {
	return &CoinService{supportedCurrencies: supportedCurrencies, supportedCoins: supportedCoins, coinClient: client, emailClient: emailClient, emailRepository: emailRepository}
}

func (c *CoinService) GetCurrentRate(currency string, coin string) (*domain.Price, error) {
	if !c.validateCurrency(currency) {
		return nil, domain.ArgumentError{Message: fmt.Sprintf("Currency %s is not supported", currency)}
	}

	if !c.validateCoin(coin) {
		return nil, domain.ArgumentError{Message: fmt.Sprintf("Coin %s is not supported", coin)}
	}

	rate, time, err := c.coinClient.GetRate(currency, coin)

	if err != nil {
		return nil, err
	}

	return &domain.Price{
		Amount:    rate,
		Currency:  currency,
		Timestamp: time,
	}, nil
}

func (c *CoinService) Subscribe(email string) error {
	err := c.emailRepository.AddEmail(email)
	if err != nil {
		return err
	}

	err = c.emailRepository.Save()
	if err != nil {
		return err
	}

	return nil
}

func (c *CoinService) SendRateEmails(currency string, coin string) error {
	if !c.validateCurrency(currency) {
		return domain.ArgumentError{Message: fmt.Sprintf("Currency %s is not supported", currency)}
	}

	if !c.validateCoin(coin) {
		return domain.ArgumentError{Message: fmt.Sprintf("Coin %s is not supported", coin)}
	}

	emails := c.emailRepository.GetAll()

	currentPrice, err := c.GetCurrentRate(currency, coin)
	if err != nil {
		return err
	}

	htmlTemplate := `<p><strong>Amount:</strong> %f</p>
	<p><strong>Currency:</strong> %s<p>
	<p><strong>Timestamp:</strong> %s<p>`
	htmlBody := fmt.Sprintf(htmlTemplate, currentPrice.Amount, currentPrice.Currency, currentPrice.Timestamp.Format("02-01-06 15:04:05.999 Z0700"))

	err = c.emailClient.Send(emails, htmlBody)
	if err != nil {
		return err
	}

	return nil
}

func (c *CoinService) validateCurrency(currency string) bool {
	return slices.Contains(c.supportedCurrencies, currency)
}

func (c *CoinService) validateCoin(coin string) bool {
	return slices.Contains(c.supportedCoins, coin)
}
