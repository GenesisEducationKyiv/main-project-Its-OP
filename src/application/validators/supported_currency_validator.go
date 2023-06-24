package validators

import "golang.org/x/exp/slices"

type SupportedCurrencyValidator struct {
	supportedCurrencies []string
}

func (v *SupportedCurrencyValidator) Validate(currency string) (bool, error) {
	return slices.Contains(v.supportedCurrencies, currency), nil
}
