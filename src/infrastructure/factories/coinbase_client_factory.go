package factories

import (
	"btcRate/application"
	"btcRate/infrastructure"
	"btcRate/infrastructure/extensions"
	"btcRate/infrastructure/integrations"
)

type CoinbaseClientFactory struct {
	logRepository extensions.ILogRepository
}

func NewCoinbaseClientFactory(logRepository extensions.ILogRepository) *CoinbaseClientFactory {
	return &CoinbaseClientFactory{logRepository: logRepository}
}

func (f *CoinbaseClientFactory) CreateClient() application.ICoinClient {
	httpClient := infrastructure.NewHttpClient(nil)
	loggedHttpClient := extensions.NewLoggedHttpClient(httpClient, f.logRepository)

	coinbaseClient := integrations.NewCoinbaseClient(loggedHttpClient)

	return coinbaseClient
}
