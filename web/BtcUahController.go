package web

import (
	"btcRate/application"
	"btcRate/docs"
	"btcRate/domain"
	"btcRate/infrastructure"
	"github.com/gin-gonic/gin"
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

const currency = "UAH"
const coin = "BTC"

var btcuahService domain.ICoinService

func RunBtcUahController() {
	var bitcoinClient = infrastructure.NewBinanceClient()
	var emailRepository = infrastructure.NewFileEmailRepository()
	var emailClient = infrastructure.NewSendGridEmailClient(os.Getenv("SENDGRID_API_KEY"), os.Getenv("SENDGRID_API_SENDER_NAME"), os.Getenv("SENDGRID_API_SENDER_EMAIL"))
	btcuahService = application.NewCoinService(bitcoinClient, emailClient, emailRepository)

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
	r.Run(":8080") // Run on port 8080
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
		c.Error(err)
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
	if email == "" {
		c.String(http.StatusBadRequest, "Email is required")
		return
	}

	if !validateEmail(&email) {
		c.String(http.StatusBadRequest, "Email is invalid")
		return
	}

	err := btcuahService.Subscribe(email)
	if err != nil {
		c.Error(err)
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
	btcuahService.SendRateEmails(currency, coin)

	c.String(http.StatusOK, "E-mails sent")
}

func validateEmail(email *string) bool {
	regexString := "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
	match, _ := regexp.Match(regexString, []byte(*email))

	return match
}
