package web

import (
	"btcRate/campaign/docs"
	"btcRate/campaign/domain"
	"btcRate/common/infrastructure"
	"btcRate/common/web"
	"context"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
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

func (*ServerManager) RunServer(fc *FileConfiguration, sc *SendgridConfiguration, pc *ProviderConfiguration, commandBus *cqrs.CommandBus) (func() error, error) {
	r := gin.Default()
	r.Use(ErrorHandlingMiddleware())

	campaignController, err := newCampaignController(fc, sc, pc, commandBus)
	if err != nil {
		return nil, err
	}

	docs.SwaggerInfo.BasePath = domain.ApiBasePath
	api := r.Group(domain.ApiBasePath)
	{
		api.POST(domain.Subscribe, campaignController.subscribe)
		api.POST(domain.SendEmails, campaignController.sendEmails)
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

func (s *ServerManager) SendEmails(host string) (*web.Response[string], error) {
	url := host + domain.ApiBasePath + domain.SendEmails

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.SendRequest(req)
	if err != nil {
		return nil, err
	}

	if isSuccessful(resp.Code) {
		result := string(resp.Body)
		if err != nil {
			return nil, err
		}
		return &web.Response[string]{Code: resp.Code, Body: &result, ErrorMessage: "", Successful: true}, nil
	}

	return &web.Response[string]{Code: resp.Code, ErrorMessage: string(resp.Body), Successful: false}, nil
}

func isSuccessful(code int) bool {
	return code >= http.StatusOK && code < http.StatusBadRequest
}
