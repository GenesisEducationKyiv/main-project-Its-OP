package factories

import (
	"coin/application/services"
	"coin/infrastructure"
	"coin/infrastructure/extensions"
	"coin/infrastructure/integrations"
)

type CoinbaseClientFactory struct {
	logRepository extensions.ILogRepository
}

func NewCoinbaseClientFactory(logRepository extensions.ILogRepository) *CoinbaseClientFactory {
	return &CoinbaseClientFactory{logRepository: logRepository}
}

func (f *CoinbaseClientFactory) CreateClient() services.ICoinClient {
	httpClient := infrastructure.NewHttpClient(nil)
	loggedHttpClient := extensions.NewLoggedHttpClient(httpClient, f.logRepository)

	coinbaseClient := integrations.NewCoinbaseClient(loggedHttpClient)

	return coinbaseClient
}
