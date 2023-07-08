package services

import "coin/application"

type IEmailRepository interface {
	AddEmail(email string) error
	GetAll() []string
}

type IEmailClient interface {
	Send(recipients []string, mailBody *application.MailBody) error
}

type CampaignService struct {
	emailRepository IEmailRepository
	emailClient     IEmailClient
	emailValidator  application.IValidator[string]
}

func NewCampaignService(repository IEmailRepository, client IEmailClient, emailValidator application.IValidator[string]) *CampaignService {
	return &CampaignService{emailRepository: repository, emailClient: client, emailValidator: emailValidator}
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

func (c *CampaignService) SendEmails(mailBody *application.MailBody) error {
	emails := c.emailRepository.GetAll()
	err := c.emailClient.Send(emails, mailBody)
	if err != nil {
		return err
	}

	return nil
}
