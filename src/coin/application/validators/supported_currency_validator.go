package validators

import (
	"btcRate/domain"
	"fmt"
	"golang.org/x/exp/slices"
)

type SupportedCurrencyValidator struct {
	supportedCurrencies []string
}

func NewSupportedCurrencyValidator(supportedCurrencies []string) *SupportedCurrencyValidator {
	return &SupportedCurrencyValidator{supportedCurrencies: supportedCurrencies}
}

func (v *SupportedCurrencyValidator) Validate(currency string) error {
	if !slices.Contains(v.supportedCurrencies, currency) {
		return &domain.ArgumentError{Message: fmt.Sprintf("currency %s is not supported", currency)}
	}

	return nil
}
