package web

import (
	"btcRate/application"
	"btcRate/application/validators"
	"btcRate/docs"
	"btcRate/domain"
	"btcRate/infrastructure"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sendgrid/sendgrid-go"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"os"
)

// @title GSES2 BTC application API
// @version 1.0.0
// @description This is a sample server for a BTC to UAH rate application.
// @host localhost:8080
// @BasePath /api

const currency = "UAH"
const coin = "BTC"

var btcuahService domain.ICoinService
var campaignService domain.ICampaignService

func RunBtcUahController(storageFile string) (func() error, error) {
	var emailRepository, err = infrastructure.NewFileEmailRepository(storageFile)
	if err != nil {
		return nil, err
	}

	var bitcoinClient = infrastructure.NewBinanceClient()
	var sendgrid = sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	var emailClient = infrastructure.NewSendGridEmailClient(sendgrid, os.Getenv("SENDGRID_API_SENDER_NAME"), os.Getenv("SENDGRID_API_SENDER_EMAIL"))

	var emailValidator = &validators.EmailValidator{}
	var supportedCoinValidator = validators.NewSupportedCoinValidator([]string{coin})
	var supportedCurrencyValidator = validators.NewSupportedCurrencyValidator([]string{currency})

	campaignService = application.NewCampaignService(emailRepository, emailClient, emailValidator)
	btcuahService = application.NewCoinService(bitcoinClient, campaignService, supportedCoinValidator, supportedCurrencyValidator)

	r := gin.Default()
	r.Use(errorHandlingMiddleware())

	docs.SwaggerInfo.BasePath = "/api/v1"
	api := r.Group("/api/v1")
	{
		api.GET("/rate", getRate)
		api.POST("/subscribe", subscribe)
		api.POST("/sendEmails", sendEmails)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	stop := func() error {
		return server.Shutdown(context.Background())
	}

	return stop, nil
}

// @Summary Get current BTC to UAH rate
// @Description Get the current rate of BTC to UAH using any third-party service with public API
// @Tags rate
// @Produce  json
// @Success 200 {number} number "Successful operation"
// @Failure 400 {object} string "Invalid status value"
// @Router /rate [get]
func getRate(c *gin.Context) {
	price, err := btcuahService.GetCurrentRate(currency, coin)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.IndentedJSON(http.StatusOK, price.Amount)
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
func subscribe(c *gin.Context) {
	email := c.PostForm("email")

	err := campaignService.Subscribe(email)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.String(http.StatusOK, "E-mail address added")
}

// @Summary Send email with BTC rate
// @Description Send the current BTC to UAH rate to all subscribed emails
// @Tags subscription
// @Produce  json
// @Success 200 {object} string "E-mails sent"
// @Router /sendEmails [post]
func sendEmails(c *gin.Context) {
	err := btcuahService.SendRateEmails(currency, coin)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.String(http.StatusOK, "E-mails sent")
}
