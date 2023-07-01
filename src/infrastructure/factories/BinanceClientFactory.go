package factories

import (
	"btcRate/infrastructure"
	"btcRate/infrastructure/integrations"
)

type BinanceClientFactory struct{}

func (*BinanceClientFactory) CreateClient() *integrations.BinanceClient {
	httpClient := infrastructure.NewExtendedHttpClient(nil)
	binanceClient := integrations.NewBinanceClient(httpClient)

	return binanceClient
}
