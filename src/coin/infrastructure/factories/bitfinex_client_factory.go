package factories

import (
	"coin/application/services"
	"coin/infrastructure"
	"coin/infrastructure/extensions"
	"coin/infrastructure/integrations"
)

type BitfinexClientFactory struct {
	logRepository extensions.ILogRepository
}

func NewBitfinexClientFactory(logRepository extensions.ILogRepository) *BitfinexClientFactory {
	return &BitfinexClientFactory{logRepository: logRepository}
}

func (f *BitfinexClientFactory) CreateClient() services.ICoinClient {
	httpClient := infrastructure.NewHttpClient(nil)
	loggedHttpClient := extensions.NewLoggedHttpClient(httpClient, f.logRepository)

	bitfinexClient := integrations.NewBitfinexClient(loggedHttpClient)

	return bitfinexClient
}
