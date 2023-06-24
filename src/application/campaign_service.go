package application

type IEmailRepository interface {
	AddEmail(email string) error
	GetAll() []string
	Save() error
}

type IEmailClient interface {
	Send(recipients []string, htmlContent string) error
}

type CampaignService struct {
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

emails := c.emailRepository.GetAll()
