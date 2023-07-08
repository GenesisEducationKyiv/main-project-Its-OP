package application

import (
	"btcRate/campaign/domain"
	"btcRate/common/application"
	"fmt"
)

type IEmailRepository interface {
	AddEmail(email string) error
	GetAll() []string
}

type IEmailClient interface {
	Send(recipients []string, mailBody *domain.MailBody) error
}

type IRateProvider interface {
	GetRate() (*domain.Rate, error)
}

type CampaignService struct {
	emailRepository IEmailRepository
	emailClient     IEmailClient
	rateProvider    IRateProvider
	emailValidator  application.IValidator[string]
}

func NewCampaignService(r IEmailRepository, c IEmailClient, rP IRateProvider, emailV application.IValidator[string]) *CampaignService {
	return &CampaignService{emailRepository: r, emailClient: c, rateProvider: rP, emailValidator: emailV}
}

func (c *CampaignService) Subscribe(email string) error {
	err := c.emailRepository.AddEmail(email)
	if err != nil {
		return err
	}

	return nil
}

func (c *CampaignService) SendRateEmails() error {
	currentPrice, err := c.rateProvider.GetRate()
	if err != nil {
		return err
	}

	htmlTemplate := `<p><strong>Amount:</strong> %f</p> <p><strong>Currency:</strong> %s<p> <p><strong>Timestamp:</strong> %s<p>`
	htmlBody := fmt.Sprintf(htmlTemplate, currentPrice.Amount, currentPrice.Currency, currentPrice.Timestamp.Format("02-01-06 15:04:05.999 Z0700"))

	mail := &domain.MailBody{Subject: "Current BTC to UAH rate", ReceiverAlias: "Rate Recipient", HtmlContent: htmlBody}

	err = c.sendEmails(mail)
	if err != nil {
		return err
	}

	return nil
}

func (c *CampaignService) sendEmails(mailBody *domain.MailBody) error {
	emails := c.emailRepository.GetAll()
	err := c.emailClient.Send(emails, mailBody)
	if err != nil {
		return err
	}

	return nil
}
