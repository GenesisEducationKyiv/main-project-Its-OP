package factories

import (
	"btcRate/application"
	"btcRate/infrastructure"
	"btcRate/infrastructure/integrations"
)

type BinanceClientFactory struct{}

func (BinanceClientFactory) CreateClient() application.ICoinClient {
	httpClient := infrastructure.NewHttpClient(nil)
	binanceClient := integrations.NewBinanceClient(httpClient)

	return binanceClient
}
