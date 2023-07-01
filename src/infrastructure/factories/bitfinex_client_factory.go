package factories

import (
	"btcRate/application"
	"btcRate/infrastructure"
	"btcRate/infrastructure/extensions"
	"btcRate/infrastructure/integrations"
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
