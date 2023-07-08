package web

import (
	"coin/application/proxies"
	"coin/application/services"
	"coin/application/validators"
	"coin/domain"
	"coin/infrastructure/factories"
	"coin/infrastructure/repositories"
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

func newBtcUahController(logStorageFile string) (*btcUahController, error) {
	supportedCurrency := "UAH"
	supportedCoin := "BTC"

	logRepository := repositories.NewFileLogRepository(logStorageFile)

	binanceFactory := factories.NewBinanceClientFactory(logRepository)
	coinbaseFactory := factories.NewCoinbaseClientFactory(logRepository)
	bitfinexFactory := factories.NewBitfinexClientFactory(logRepository)
	coinClientFactories := []services.ICoinClientFactory{binanceFactory, coinbaseFactory, bitfinexFactory}

	chainedCoinClientFactory := proxies.NewChainedCoinClientFactory(coinClientFactories)

	var supportedCoinValidator = validators.NewSupportedCoinValidator([]string{supportedCoin})
	var supportedCurrencyValidator = validators.NewSupportedCurrencyValidator([]string{supportedCurrency})

	var btcUahService = services.NewCoinService(chainedCoinClientFactory, supportedCoinValidator, supportedCurrencyValidator)

	controller := &btcUahController{coinService: btcUahService, currency: supportedCurrency, coin: supportedCoin}

	return controller, nil
}

// @Summary Get current BTC to UAH rate
// @Description Get the current rate of BTC to UAH using any third-party coinService with public API
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

	context.IndentedJSON(http.StatusOK, price.Amount)
}
