package factories

import (
	"btcRate/coin/application"
	"btcRate/coin/infrastructure"
	"btcRate/coin/infrastructure/extensions"
	"btcRate/coin/infrastructure/integrations"
)

type BinanceClientFactory struct {
	logRepository extensions.ILogRepository
}

func NewBinanceClientFactory(logRepository extensions.ILogRepository) *BinanceClientFactory {
	return &BinanceClientFactory{logRepository: logRepository}
}

func (f *BinanceClientFactory) CreateClient() application.ICoinClient {
	httpClient := infrastructure.NewHttpClient(nil)
	loggedHttpClient := extensions.NewLoggedHttpClient(httpClient, f.logRepository)

	binanceClient := integrations.NewBinanceClient(loggedHttpClient)

	return binanceClient
}
