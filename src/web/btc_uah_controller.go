package web

import (
	"btcRate/application"
	"btcRate/docs"
	"btcRate/domain"
	"btcRate/infrastructure"
	"github.com/gin-gonic/gin"
	"github.com/sendgrid/sendgrid-go"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

func RunBtcUahController() error {
	var emailRepository, err = infrastructure.NewFileEmailRepository()
	if err != nil {
		return err
	}

	var bitcoinClient = infrastructure.NewBinanceClient()
	var sendgrid = sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	var emailClient = infrastructure.NewSendGridEmailClient(sendgrid, os.Getenv("SENDGRID_API_SENDER_NAME"), os.Getenv("SENDGRID_API_SENDER_EMAIL"))
	var btcUahService = application.NewCoinService(bitcoinClient, emailClient, emailRepository)

	controller := btcUahController{service: btcUahService, currency: "UAH", coin: "BTC"}

	r := gin.Default()
	r.Use(errorHandlingMiddleware())

	docs.SwaggerInfo.BasePath = "/api/v1"
	api := r.Group("/api/v1")
	{
		api.GET("/rate", controller.getRate)
		api.POST("/subscribe", controller.subscribe)
		api.POST("/sendEmails", controller.sendEmails)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	err = r.Run(":8080") // Run on port 8080
	if err != nil {
		return err
	}

	return nil
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
