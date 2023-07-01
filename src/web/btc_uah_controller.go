package web

import (
	"btcRate/application"
	"btcRate/application/validators"
	"btcRate/domain"
	"btcRate/infrastructure/factories"
	"btcRate/infrastructure/integrations"
	"btcRate/infrastructure/repositories"
	"github.com/gin-gonic/gin"
	"github.com/sendgrid/sendgrid-go"
	"net/http"
	"os"
)

// @title GSES2 BTC application API
// @version 1.0.0
// @description This is a sample server for a BTC to UAH rate application.
// @host localhost:8080
// @BasePath /api

type btcUahController struct {
	coinService     domain.ICoinService
	campaignService domain.ICampaignService
	currency        string
	coin            string
}

func newBtcUahController(storageFile string) (*btcUahController, error) {
	supportedCurrency := "UAH"
	supportedCoin := "BTC"

	var emailRepository, err = repositories.NewFileEmailRepository(storageFile)
	if err != nil {
		return nil, err
	}

	binanceFactory := factories.BinanceClientFactory{}
	coinbaseFactory := factories.CoinbaseClientFactory{}
	bitfinexFactory := factories.BitfinexClientFactory{}

	coinClientFactories := []application.ICoinClientFactory{binanceFactory, coinbaseFactory, bitfinexFactory}

	var sendgrid = sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	var emailClient = integrations.NewSendGridEmailClient(sendgrid, os.Getenv("SENDGRID_API_SENDER_NAME"), os.Getenv("SENDGRID_API_SENDER_EMAIL"))

	var emailValidator = &validators.EmailValidator{}
	var supportedCoinValidator = validators.NewSupportedCoinValidator([]string{supportedCoin})
	var supportedCurrencyValidator = validators.NewSupportedCurrencyValidator([]string{supportedCurrency})

	var campaignService = application.NewCampaignService(emailRepository, emailClient, emailValidator)

	var btcUahService = application.NewCoinService(coinClientFactories, campaignService, supportedCoinValidator, supportedCurrencyValidator)

	controller := &btcUahController{coinService: btcUahService, campaignService: campaignService, currency: supportedCurrency, coin: supportedCoin}

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

// @Summary Subscribe email to get BTC rate
// @Description Add an email to the database if it does not exist already
// @Tags subscription
// @Accept  x-www-form-urlencoded
// @Produce  json
// @Param email formData string true "Email to be subscribed"
// @Success 200 {object} string "E-mail added"
// @Failure 409 {object} string "E-mail already exists in the database"
// @Router /subscribe [post]
func (c *btcUahController) subscribe(context *gin.Context) {
	email := context.PostForm("email")

	err := c.campaignService.Subscribe(email)
	if err != nil {
		_ = context.Error(err)
		return
	}

	context.String(http.StatusOK, "E-mail address added")
}

// @Summary Send email with BTC rate
// @Description Send the current BTC to UAH rate to all subscribed emails
// @Tags subscription
// @Produce  json
// @Success 200 {object} string "E-mails sent"
// @Router /sendEmails [post]
func (c *btcUahController) sendEmails(context *gin.Context) {
	err := c.coinService.SendRateEmails(c.currency, c.coin)
	if err != nil {
		_ = context.Error(err)
		return
	}

	context.String(http.StatusOK, "E-mails sent")
}
