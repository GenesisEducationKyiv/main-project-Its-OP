package application

import (
	"btcRate/domain"
	"fmt"
	"time"
)

type ICoinClientFactory interface {
	CreateClient() ICoinClient
}

type ICoinClient interface {
	GetRate(currency string, coin string) (float64, time.Time, error)
}

type CoinService struct {
	coinClient        ICoinClient
	campaignService   domain.ICampaignService
	coinValidator     domain.IValidator[string]
	currencyValidator domain.IValidator[string]
}

func NewCoinService(client ICoinClient, campaignService domain.ICampaignService, coinValidator domain.IValidator[string], currencyValidator domain.IValidator[string]) *CoinService {
	return &CoinService{coinClient: client, campaignService: campaignService, coinValidator: coinValidator, currencyValidator: currencyValidator}
}

func (c *CoinService) GetCurrentRate(currency string, coin string) (*domain.Price, error) {
	err := c.validateParameters(currency, coin)
	if err != nil {
		return nil, err
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

func (c *CoinService) SendRateEmails(currency string, coin string) error {
	err := c.validateParameters(currency, coin)
	if err != nil {
		return err
	}

	currentPrice, err := c.GetCurrentRate(currency, coin)
	if err != nil {
		return err
	}

	htmlTemplate := `<p><strong>Amount:</strong> %f</p>
	<p><strong>Currency:</strong> %s<p>
	<p><strong>Timestamp:</strong> %s<p>`
	htmlBody := fmt.Sprintf(htmlTemplate, currentPrice.Amount, currentPrice.Currency, currentPrice.Timestamp.Format("02-01-06 15:04:05.999 Z0700"))

	err = c.campaignService.SendEmails(htmlBody)
	if err != nil {
		return err
	}

	return nil
}

func (c *CoinService) validateParameters(currency string, coin string) error {
	err := c.currencyValidator.Validate(currency)
	if err != nil {
		return err
	}

	err = c.coinValidator.Validate(coin)
	if err != nil {
		return err
	}

	return nil
}
