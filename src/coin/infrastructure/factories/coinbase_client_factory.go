package factories

import (
	"btcRate/coin/application"
	"btcRate/coin/infrastructure/integrations"
	commonApplication "btcRate/common/application"
	"btcRate/common/infrastructure"
	"btcRate/common/infrastructure/extensions"
)

type CoinbaseClientFactory struct {
	logger commonApplication.ILogger
}

func NewCoinbaseClientFactory(logger commonApplication.ILogger) *CoinbaseClientFactory {
	return &CoinbaseClientFactory{logger: logger}
}

func (f *CoinbaseClientFactory) CreateClient() application.ICoinClient {
	httpClient := infrastructure.NewHttpClient(nil)
	loggedHttpClient := extensions.NewLoggedHttpClient(httpClient, f.logger)

	coinbaseClient := integrations.NewCoinbaseClient(loggedHttpClient)

	return coinbaseClient
}
