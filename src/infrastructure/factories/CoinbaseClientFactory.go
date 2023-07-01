package factories

import (
	"btcRate/application"
	"btcRate/infrastructure"
	"btcRate/infrastructure/integrations"
)

type CoinbaseClientFactory struct{}

func (CoinbaseClientFactory) CreateClient() application.ICoinClient {
	httpClient := infrastructure.NewHttpClient(nil)
	coinbaseClient := integrations.NewCoinbaseClient(httpClient)

	return coinbaseClient
}
