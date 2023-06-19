package infrastructure

import (
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"log"
)

type SendGridEmailClient struct {
	client      *sendgrid.Client
	senderName  string
	senderEmail string
}

func NewSendGridEmailClient(client *sendgrid.Client, senderName string, senderEmail string) *SendGridEmailClient {
	return &SendGridEmailClient{client: client, senderName: senderName, senderEmail: senderEmail}
}

func (s *SendGridEmailClient) Send(recipients []string, htmlContent string) {
	if len(recipients) == 0 {
		return
	}

	from := mail.NewEmail(s.senderName, s.senderEmail)
	firstTo := mail.NewEmail("Rate Recipient", recipients[0])
	subject := "Current BTC to UAH rate"
	message := mail.NewSingleEmail(from, subject, firstTo, "", htmlContent)

	for i := 1; i < len(recipients); i++ {
		personalization := mail.NewPersonalization()
		personalization.AddTos(mail.NewEmail("Rate Recipient", recipients[i]))
		message.AddPersonalizations(personalization)
	}

	response, err := s.client.Send(message)
	if err != nil {
		log.Fatalln(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}
