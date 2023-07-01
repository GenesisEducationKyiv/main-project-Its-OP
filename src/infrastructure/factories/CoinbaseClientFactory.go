package factories

import (
	"btcRate/infrastructure"
	"btcRate/infrastructure/integrations"
)

type CoinbaseClientFactory struct{}

func (*CoinbaseClientFactory) CreateClient() *integrations.CoinbaseClient {
	httpClient := infrastructure.NewExtendedHttpClient(nil)
	coinbaseClient := integrations.NewCoinbaseClient(httpClient)

	return coinbaseClient
}
