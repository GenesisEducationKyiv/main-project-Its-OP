package validators

import "golang.org/x/exp/slices"

type SupportedCoinValidator struct {
	supportedCoins []string
}

func NewSupportedCoinValidator(supportedCoins []string) *SupportedCoinValidator {
	return &SupportedCoinValidator{supportedCoins: supportedCoins}
}

func (v *SupportedCoinValidator) Validate(coin string) (bool, error) {
	return slices.Contains(v.supportedCoins, coin), nil
}
