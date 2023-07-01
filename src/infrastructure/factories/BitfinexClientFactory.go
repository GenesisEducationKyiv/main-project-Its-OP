package factories

import (
	"btcRate/application"
	"btcRate/infrastructure"
	"btcRate/infrastructure/integrations"
)

type BitfinexClientFactory struct{}

func (BitfinexClientFactory) CreateClient() application.ICoinClient {
	httpClient := infrastructure.NewHttpClient(nil)
	bitfinexClient := integrations.NewBitfinexClient(httpClient)

	return bitfinexClient
}
