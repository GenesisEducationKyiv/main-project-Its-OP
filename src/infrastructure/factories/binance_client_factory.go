package factories

import (
	"btcRate/application/services"
	"btcRate/infrastructure"
	"btcRate/infrastructure/extensions"
	"btcRate/infrastructure/integrations"
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
