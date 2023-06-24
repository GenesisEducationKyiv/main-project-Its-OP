package validators

import "regexp"

type EmailValidator struct{}

func (e *EmailValidator) Validate(email string) (bool, error) {
	regexString := "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"

	match, err := regexp.Match(regexString, []byte(email))
	if err != nil {
		return false, err
	}

	return match, nil
}
