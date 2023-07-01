package factories

import (
	"btcRate/application"
	"btcRate/infrastructure"
	"btcRate/infrastructure/extensions"
	"btcRate/infrastructure/integrations"
)

type BinanceClientFactory struct{}

func (BinanceClientFactory) CreateClient() application.ICoinClient {
	httpClient := infrastructure.NewHttpClient(nil)
	loggedHttpClient := extensions.NewLoggedHttpClient(httpClient)

	binanceClient := integrations.NewBinanceClient(loggedHttpClient)

	return binanceClient
}
