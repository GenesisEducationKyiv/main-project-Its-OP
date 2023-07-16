package web

import (
	"btcRate/coin/application"
	"btcRate/coin/application/proxies"
	"btcRate/coin/application/validators"
	"btcRate/coin/domain"
	"btcRate/coin/infrastructure/factories"
	"btcRate/common/infrastructure"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @title GSES2 BTC application API
// @version 1.0.0
// @description This is a sample server for a BTC to UAH rate application.
// @host localhost:8080
// @BasePath /api

type btcUahController struct {
	coinService domain.ICoinService
	currency    string
	coin        string
}

func newBtcUahController(logStorageFile string, commandBus *cqrs.CommandBus) (*btcUahController, error) {
	supportedCurrency := "UAH"
	supportedCoin := "BTC"

	logger := infrastructure.NewLogger(commandBus)

	binanceFactory := factories.NewBinanceClientFactory(logger)
	coinbaseFactory := factories.NewCoinbaseClientFactory(logger)
	bitfinexFactory := factories.NewBitfinexClientFactory(logger)
	coinClientFactories := []application.ICoinClientFactory{binanceFactory, coinbaseFactory, bitfinexFactory}

	chainedCoinClientFactory := proxies.NewChainedCoinClientFactory(coinClientFactories)

	var supportedCoinValidator = validators.NewSupportedCoinValidator([]string{supportedCoin})
	var supportedCurrencyValidator = validators.NewSupportedCurrencyValidator([]string{supportedCurrency})

	var btcUahService = application.NewCoinService(chainedCoinClientFactory, supportedCoinValidator, supportedCurrencyValidator, logger)

	controller := &btcUahController{coinService: btcUahService, currency: supportedCurrency, coin: supportedCoin}

	return controller, nil
}

// @Summary Get current BTC to UAH rate
// @Description Get the current rate of BTC to UAH using any third-party service with public API
// @Tags rate
// @Produce  json
// @Success 200 {number} number "Successful operation"
// @Failure 400 {object} string "Invalid status value"
// @Router /rate [get]
func (c *btcUahController) getRate(context *gin.Context) {
	price, err := c.coinService.GetCurrentRate(c.currency, c.coin)

	if err != nil {
		_ = context.Error(err)
		return
	}

	context.IndentedJSON(http.StatusOK, price)
}
