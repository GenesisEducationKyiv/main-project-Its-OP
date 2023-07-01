package web

import (
	"btcRate/application"
	"btcRate/domain"
	"btcRate/infrastructure"
	"github.com/gin-gonic/gin"
	"github.com/sendgrid/sendgrid-go"
	"net/http"
	"os"
	"regexp"
)

// @title GSES2 BTC application API
// @version 1.0.0
// @description This is a sample server for a BTC to UAH rate application.
// @host localhost:8080
// @BasePath /api

type btcUahController struct {
	service  domain.ICoinService
	currency string
	coin     string
}

func newBtcUahController(storageFile string) (*btcUahController, error) {
	var emailRepository, err = infrastructure.NewFileEmailRepository(storageFile)
	if err != nil {
		return nil, err
	}

	var bitcoinClient = infrastructure.NewBinanceClient()
	var sendgrid = sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	var emailClient = infrastructure.NewSendGridEmailClient(sendgrid, os.Getenv("SENDGRID_API_SENDER_NAME"), os.Getenv("SENDGRID_API_SENDER_EMAIL"))

	supportedCurrency := "UAH"
	supportedCoin := "BTC"
	var btcUahService = application.NewCoinService([]string{supportedCurrency}, []string{supportedCoin}, bitcoinClient, emailClient, emailRepository)

	controller := &btcUahController{service: btcUahService, currency: supportedCurrency, coin: supportedCoin}

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
	price, err := c.service.GetCurrentRate(c.currency, c.coin)

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
	if email == "" {
		context.String(http.StatusBadRequest, "Email is required")
		return
	}

	if valid, err := validateEmail(&email); err != nil {
		_ = context.Error(err)
		return
	} else if !valid {
		context.String(http.StatusBadRequest, "Email is invalid")
		return
	}

	err := c.service.Subscribe(email)
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
	err := c.service.SendRateEmails(c.currency, c.coin)
	if err != nil {
		_ = context.Error(err)
		return
	}

	context.String(http.StatusOK, "E-mails sent")
}

func validateEmail(email *string) (bool, error) {
	regexString := "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
	match, err := regexp.Match(regexString, []byte(*email))
	if err != nil {
		return false, err
	}

	return match, nil
}
