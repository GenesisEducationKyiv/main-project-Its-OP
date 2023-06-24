package validators

import (
	"fmt"
	"regexp"
)

type EmailValidator struct{}

func (v *EmailValidator) Validate(email string) error {
	regexString := "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"

	match, err := regexp.Match(regexString, []byte(email))
	if err != nil {
		return err
	}
	if !match {
		return fmt.Errorf("email is invalid")
	}

	return nil
}
