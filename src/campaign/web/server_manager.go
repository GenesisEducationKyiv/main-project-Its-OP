package web

import (
	"btcRate/campaign/docs"
	"btcRate/common/infrastructure"
	"btcRate/common/web"
	"context"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
)

type ServerManager struct {
	client infrastructure.IHttpClient
}

type Response[T any] struct {
	Code         int
	Body         *T
	ErrorMessage string
	Successful   bool
}

func NewServerManager() ServerManager {
	return ServerManager{infrastructure.NewHttpClient(nil)}
}

func (*ServerManager) RunServer(fc *FileConfiguration, sc *SendgridConfiguration, pc *ProviderConfiguration) (func() error, error) {
	r := gin.Default()
	r.Use(web.ErrorHandlingMiddleware())

	campaignController, err := newCampaignController(fc, sc, pc)
	if err != nil {
		return nil, err
	}

	docs.SwaggerInfo.BasePath = web.ApiBasePath
	api := r.Group(web.ApiBasePath)
	{
		api.POST(web.Subscribe, campaignController.subscribe)
		api.POST(web.SendEmails, campaignController.sendEmails)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	server := &http.Server{
		Addr:    ":8081",
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
