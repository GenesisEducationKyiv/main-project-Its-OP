package factories

import (
	"coin/application/services"
	"coin/infrastructure"
	"coin/infrastructure/extensions"
	"coin/infrastructure/integrations"
)

type BinanceClientFactory struct {
	logRepository extensions.ILogRepository
}

func NewBinanceClientFactory(logRepository extensions.ILogRepository) *BinanceClientFactory {
	return &BinanceClientFactory{logRepository: logRepository}
}

func (f *BinanceClientFactory) CreateClient() services.ICoinClient {
	httpClient := infrastructure.NewHttpClient(nil)
	loggedHttpClient := extensions.NewLoggedHttpClient(httpClient, f.logRepository)

	binanceClient := integrations.NewBinanceClient(loggedHttpClient)

	return binanceClient
}
