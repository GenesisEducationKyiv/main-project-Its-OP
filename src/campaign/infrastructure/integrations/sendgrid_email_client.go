package integrations

import (
	"btcRate/campaign/domain"
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"net/http"
)

type SendGridEmailClient struct {
	client      *sendgrid.Client
	senderName  string
	senderEmail string
}

func NewSendGridEmailClient(client *sendgrid.Client, senderName string, senderEmail string) *SendGridEmailClient {
	return &SendGridEmailClient{client: client, senderName: senderName, senderEmail: senderEmail}
}

func (s *SendGridEmailClient) Send(recipients []string, mailBody *domain.MailBody) error {
	if len(recipients) == 0 {
		return nil
	}

	from := mail.NewEmail(s.senderName, s.senderEmail)
	firstTo := mail.NewEmail(mailBody.ReceiverAlias, recipients[0])
	subject := mailBody.Subject
	message := mail.NewSingleEmail(from, subject, firstTo, "", mailBody.HtmlContent)

	for i := 1; i < len(recipients); i++ {
		personalization := mail.NewPersonalization()
		personalization.AddTos(mail.NewEmail(mailBody.ReceiverAlias, recipients[i]))
		message.AddPersonalizations(personalization)
	}

	response, err := s.client.Send(message)
	if err != nil {
		return err
	} else if response.StatusCode < http.StatusOK || response.StatusCode > http.StatusNoContent {
		err = fmt.Errorf("error sending an email. %d: %s", response.StatusCode, response.Body)
		return err
	}

	return nil
}
