package infrastructure

import (
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"log"
)

type SendGridEmailClient struct {
	apiKey      string
	senderName  string
	senderEmail string
}

func NewSendGridEmailClient(apiKey string, senderName string, senderEmail string) *SendGridEmailClient {
	return &SendGridEmailClient{apiKey: apiKey, senderName: senderName, senderEmail: senderEmail}
}

func (emailClient *SendGridEmailClient) Send(recipients []string, htmlContent string) {
	if len(recipients) == 0 {
		return
	}

	from := mail.NewEmail(emailClient.senderName, emailClient.senderEmail)
	firstTo := mail.NewEmail("Rate Recipient", recipients[0])
	subject := "Current BTC to UAH rate"
	message := mail.NewSingleEmail(from, subject, firstTo, "", htmlContent)

	for i := 1; i < len(recipients); i++ {
		personalization := mail.NewPersonalization()
		personalization.AddTos(mail.NewEmail("Rate Recipient", recipients[i]))
		message.AddPersonalizations(personalization)
	}

	client := sendgrid.NewSendClient(emailClient.apiKey)

	response, err := client.Send(message)
	if err != nil {
		log.Fatalln(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}
