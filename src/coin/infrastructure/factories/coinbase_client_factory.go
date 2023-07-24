package factories

import (
	"btcRate/coin/application"
	"btcRate/coin/infrastructure/integrations"
	"btcRate/common/infrastructure"
	"btcRate/common/infrastructure/extensions"
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
