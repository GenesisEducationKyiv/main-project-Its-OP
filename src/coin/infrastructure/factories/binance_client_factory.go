package factories

import (
	"btcRate/coin/application"
	"btcRate/coin/infrastructure/integrations"
	commonApplication "btcRate/common/application"
	"btcRate/common/infrastructure"
	"btcRate/common/infrastructure/extensions"
)

type BinanceClientFactory struct {
	logger commonApplication.ILogger
}

func NewBinanceClientFactory(logger commonApplication.ILogger) *BinanceClientFactory {
	return &BinanceClientFactory{logger: logger}
}

func (f *BinanceClientFactory) CreateClient() application.ICoinClient {
	httpClient := infrastructure.NewHttpClient(nil)
	loggedHttpClient := extensions.NewLoggedHttpClient(httpClient, f.logger)

	binanceClient := integrations.NewBinanceClient(loggedHttpClient)

	return binanceClient
}
