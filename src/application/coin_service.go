package application

import (
	"btcRate/domain"
	"fmt"
	"time"
)

//go:generate mockery --name ICoinClient
type ICoinClient interface {
	GetRate(currency string, coin string) (*SpotPrice, error)
}

type CoinService struct {
	coinClient        ICoinClient
	campaignService   domain.ICampaignService
	coinValidator     domain.IValidator[string]
	currencyValidator domain.IValidator[string]
}

type SpotPrice struct {
	Amount    float64
	Timestamp time.Time
}

func NewCoinService(client ICoinClient, campaignService domain.ICampaignService, coinValidator domain.IValidator[string], currencyValidator domain.IValidator[string]) *CoinService {
	return &CoinService{coinClient: client, campaignService: campaignService, coinValidator: coinValidator, currencyValidator: currencyValidator}
}

func (c *CoinService) GetCurrentRate(currency string, coin string) (*domain.Price, error) {
	err := c.validateParameters(currency, coin)
	if err != nil {
		return nil, err
	}

	price, err := c.coinClient.GetRate(currency, coin)

	if err != nil {
		return nil, err
	}

	return &domain.Price{
		Amount:    price.Amount,
		Currency:  currency,
		Timestamp: price.Timestamp,
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
