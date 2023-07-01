package factories

import (
	"btcRate/application"
	"btcRate/infrastructure"
	"btcRate/infrastructure/extensions"
	"btcRate/infrastructure/integrations"
)

type CoinbaseClientFactory struct{}

func (CoinbaseClientFactory) CreateClient() application.ICoinClient {
	httpClient := infrastructure.NewHttpClient(nil)
	loggedHttpClient := extensions.NewLoggedHttpClient(httpClient)

	coinbaseClient := integrations.NewCoinbaseClient(loggedHttpClient)

	return coinbaseClient
}
