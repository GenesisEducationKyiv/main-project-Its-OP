package factories

import (
	"btcRate/coin/application"
	"btcRate/coin/infrastructure/integrations"
	"btcRate/common/infrastructure"
	"btcRate/common/infrastructure/extensions"
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
