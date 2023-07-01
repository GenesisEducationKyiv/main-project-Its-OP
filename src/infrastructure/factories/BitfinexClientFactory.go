package factories

import (
	"btcRate/application"
	"btcRate/infrastructure"
	"btcRate/infrastructure/extensions"
	"btcRate/infrastructure/integrations"
)

type BitfinexClientFactory struct{}

func (BitfinexClientFactory) CreateClient() application.ICoinClient {
	httpClient := infrastructure.NewHttpClient(nil)
	loggedHttpClient := extensions.NewLoggedHttpClient(httpClient)

	bitfinexClient := integrations.NewBitfinexClient(loggedHttpClient)

	return bitfinexClient
}
