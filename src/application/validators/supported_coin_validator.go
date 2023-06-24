package validators

import (
	"btcRate/domain"
	"fmt"
	"golang.org/x/exp/slices"
)

type SupportedCoinValidator struct {
	supportedCoins []string
}

func NewSupportedCoinValidator(supportedCoins []string) *SupportedCoinValidator {
	return &SupportedCoinValidator{supportedCoins: supportedCoins}
}

func (v *SupportedCoinValidator) Validate(coin string) error {
	if slices.Contains(v.supportedCoins, coin) {
		return domain.ArgumentError{Message: fmt.Sprintf("coin %s is not supported", coin)}
	}

	return nil
}
