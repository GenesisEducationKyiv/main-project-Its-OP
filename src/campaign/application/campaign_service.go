package application

import (
	"campaign/domain"
)

type IEmailRepository interface {
	AddEmail(email string) error
	GetAll() []string
}

type IEmailClient interface {
	Send(recipients []string, mailBody *domain.MailBody) error
}

type IValidator[T any] interface {
	Validate(T) error
}

type IRateProvider interface {
	GetRate() (domain.Rate, error)
}

type CampaignService struct {
	emailRepository IEmailRepository
	emailClient     IEmailClient
	rateProvider    IRateProvider
	emailValidator  IValidator[string]
}

func NewCampaignService(repository IEmailRepository, client IEmailClient, rateProvider IRateProvider, emailValidator IValidator[string]) *CampaignService {
	return &CampaignService{emailRepository: repository, emailClient: client, rateProvider: rateProvider, emailValidator: emailValidator}
}

func (c *CampaignService) Subscribe(email string) error {
	err := c.emailValidator.Validate(email)
	if err != nil {
		return err
	}

	err = c.emailRepository.AddEmail(email)
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
