package application

import "btcRate/domain"

type IEmailRepository interface {
	AddEmail(email string) error
	GetAll() []string
	Save() error
}

type IEmailClient interface {
	Send(recipients []string, htmlContent string) error
}

type CampaignService struct {
	emailRepository IEmailRepository
	emailClient     IEmailClient
	emailValidator  domain.IValidator[string]
}

func NewCampaignService(repository IEmailRepository, client IEmailClient, emailValidator domain.IValidator[string]) *CampaignService {
	return &CampaignService{emailRepository: repository, emailClient: client, emailValidator: emailValidator}
}

func (c *CampaignService) Subscribe(email string) error {
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

func (c *CampaignService) SendEmails(htmlBody string) error {
	emails := c.emailRepository.GetAll()
	err := c.emailClient.Send(emails, htmlBody)
	if err != nil {
		return err
	}

	return nil
}
