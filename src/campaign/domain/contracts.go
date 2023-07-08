package domain

type ICampaignService interface {
	Subscribe(email string) error
	SendEmails(body *MailBody) error
}

type IValidator[T any] interface {
	Validate(T) error
}
