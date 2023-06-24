package validators

import "golang.org/x/exp/slices"

type SupportedCoinValidator struct {
	supportedCoins []string
}

func (v *SupportedCoinValidator) Validate(coin string) (bool, error) {
	return slices.Contains(v.supportedCoins, coin), nil
}
