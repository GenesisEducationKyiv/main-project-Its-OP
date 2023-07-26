package factories

import (
	"btcRate/coin/application"
	"btcRate/coin/infrastructure/integrations"
	commonApplication "btcRate/common/application"
	"btcRate/common/infrastructure"
	"btcRate/common/infrastructure/extensions"
)

type BitfinexClientFactory struct {
	logger commonApplication.ILogger
}

func NewBitfinexClientFactory(logger commonApplication.ILogger) *BitfinexClientFactory {
	return &BitfinexClientFactory{logger: logger}
}

func (f *BitfinexClientFactory) CreateClient() application.ICoinClient {
	httpClient := infrastructure.NewHttpClient(nil)
	loggedHttpClient := extensions.NewLoggedHttpClient(httpClient, f.logger)

	bitfinexClient := integrations.NewBitfinexClient(loggedHttpClient)

	return bitfinexClient
}
