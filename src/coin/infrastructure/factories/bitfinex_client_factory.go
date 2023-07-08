package factories

import (
	"btcRate/coin/application"
	"btcRate/coin/infrastructure"
	"btcRate/coin/infrastructure/extensions"
	"btcRate/coin/infrastructure/integrations"
)

type BitfinexClientFactory struct {
	logRepository extensions.ILogRepository
}

func NewBitfinexClientFactory(logRepository extensions.ILogRepository) *BitfinexClientFactory {
	return &BitfinexClientFactory{logRepository: logRepository}
}

func (f *BitfinexClientFactory) CreateClient() application.ICoinClient {
	httpClient := infrastructure.NewHttpClient(nil)
	loggedHttpClient := extensions.NewLoggedHttpClient(httpClient, f.logRepository)

	bitfinexClient := integrations.NewBitfinexClient(loggedHttpClient)

	return bitfinexClient
}
