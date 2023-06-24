package validators

import (
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
	if slices.Contains(v.supportedCurrencies, currency) {
		return fmt.Errorf("currency %s is not supported", currency)
	}

	return nil
}
