package web

import (
	"btcRate/campaign/application"
	"btcRate/campaign/application/validators"
	"btcRate/campaign/infrastructure/integrations"
	"btcRate/campaign/infrastructure/providers"
	"btcRate/campaign/infrastructure/repositories"
	"btcRate/common/infrastructure"
	"btcRate/common/infrastructure/extensions"
	commonRepositories "btcRate/common/infrastructure/repositories"
	"btcRate/common/web"
	"github.com/gin-gonic/gin"
	"github.com/sendgrid/sendgrid-go"
	"net/http"
	"net/url"
	"sync"
)

type ICampaignService interface {
	Subscribe(email string) error
	SendRateEmails() error
}

// @title GSES2 BTC application API
// @version 1.0.0
// @description This is a sample server for a BTC to UAH rate application.
// @host localhost:8081
// @BasePath /api

type campaignController struct {
	campaignService ICampaignService
	currency        string
	coin            string
}

func newCampaignController(fc *FileConfiguration, sc *SendgridConfiguration, pc *ProviderConfiguration) (*campaignController, error) {
	supportedCurrency := "UAH"
	supportedCoin := "BTC"

	emailMutex := &sync.RWMutex{}

	var emailRepository, err = repositories.NewFileEmailRepository(fc.EmailStorageFile, emailMutex)
	if err != nil {
		return nil, err
	}

	logRepository, err := commonRepositories.NewFileLogRepository(fc.LogStorageFile)
	if err != nil {
		return nil, err
	}

	var sendgrid = sendgrid.NewSendClient(sc.ApiKey)
	var emailClient = integrations.NewSendGridEmailClient(sendgrid, sc.SenderName, sc.SenderEmail)

	httpClient := infrastructure.NewHttpClient(nil)
	loggedHttpClient := extensions.NewLoggedHttpClient(httpClient, logRepository)

	var emailValidator = &validators.EmailValidator{}

	var rateProvider = providers.NewRateProvider(loggedHttpClient, &url.URL{Scheme: pc.Schema, Host: pc.Hostname, Path: web.ApiBasePath})

	var campaignService = application.NewCampaignService(emailRepository, emailClient, rateProvider, emailValidator)

	controller := &campaignController{campaignService: campaignService, currency: supportedCurrency, coin: supportedCoin}

	return controller, nil
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
func (c *campaignController) subscribe(context *gin.Context) {
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
func (c *campaignController) sendEmails(context *gin.Context) {
	err := c.campaignService.SendRateEmails()
	if err != nil {
		_ = context.Error(err)
		return
	}

	context.String(http.StatusOK, "E-mails sent")
}
