package web

import (
	"btcRate/coin/docs"
	"btcRate/coin/domain"
	"btcRate/common/infrastructure"
	"btcRate/common/web"
	"context"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
)

type ServerManager struct {
	client     infrastructure.IHttpClient
	commandBus *cqrs.CommandBus
}

func NewServerManager(commandBus *cqrs.CommandBus) ServerManager {
	return ServerManager{infrastructure.NewHttpClient(nil), commandBus}
}

func (s *ServerManager) RunServer(logStorageFile string) (func() error, error) {
	r := gin.Default()
	r.Use(web.ErrorHandlingMiddleware())

	btcUahController, err := newBtcUahController(logStorageFile, s.commandBus)
	if err != nil {
		return nil, err
	}

	docs.SwaggerInfo.BasePath = web.ApiBasePath
	api := r.Group(web.ApiBasePath)
	{
		api.GET(web.GetRate, btcUahController.getRate)
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

func (s *ServerManager) GetRate(host string) (*web.Response[domain.Price], error) {
	url := host + web.ApiBasePath + web.GetRate

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.SendRequest(req)
	if err != nil {
		return nil, err
	}

	if isSuccessful(resp.Code) {
		var result domain.Price
		err = json.Unmarshal(resp.Body, &result)
		if err != nil {
			return nil, err
		}
		return &web.Response[domain.Price]{Code: resp.Code, Body: &result, ErrorMessage: "", Successful: true}, nil
	}

	return &web.Response[domain.Price]{Code: resp.Code, ErrorMessage: string(resp.Body), Successful: false}, nil
}

func isSuccessful(code int) bool {
	return code >= http.StatusOK && code < http.StatusBadRequest
}
