package factories

import (
	"btcRate/infrastructure"
	"btcRate/infrastructure/integrations"
)

type BitfinexClientFactory struct{}

func (*BitfinexClientFactory) CreateClient() *integrations.BitfinexClient {
	httpClient := infrastructure.NewExtendedHttpClient(nil)
	binanceClient := integrations.NewBitfinexClient(httpClient)

	return binanceClient
}
