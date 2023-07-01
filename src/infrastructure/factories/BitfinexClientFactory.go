package factories

import (
	"btcRate/application"
	"btcRate/infrastructure"
	"btcRate/infrastructure/integrations"
)

type BitfinexClientFactory struct{}

func (BitfinexClientFactory) CreateClient() application.ICoinClient {
	httpClient := infrastructure.NewExtendedHttpClient(nil)
	bitfinexClient := integrations.NewBitfinexClient(httpClient)

	return bitfinexClient
}
