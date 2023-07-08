package web

import (
	"campaign/docs"
	"campaign/infrastructure"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
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

func (*ServerManager) RunServer(emailStorageFile string, logStorageFile string) (func() error, error) {
	r := gin.Default()
	r.Use(errorHandlingMiddleware())

	btcUahController, err := newCampaignController(emailStorageFile, logStorageFile)
	if err != nil {
		return nil, err
	}

	docs.SwaggerInfo.BasePath = apiBasePath
	api := r.Group(apiBasePath)
	{
		api.POST(subscribe, btcUahController.subscribe)
		api.POST(sendEmails, btcUahController.sendEmails)
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

func (s *ServerManager) GetRate(host string) (*Response[int], error) {
	url := host + apiBasePath + getRate

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.SendRequest(req)
	if err != nil {
		return nil, err
	}

	if isSuccessful(resp.Code) {
		var result int
		err = json.Unmarshal(resp.Body, &result)
		if err != nil {
			return nil, err
		}
		return &Response[int]{Code: resp.Code, Body: &result, ErrorMessage: "", Successful: true}, nil
	}

	return &Response[int]{Code: resp.Code, ErrorMessage: string(resp.Body), Successful: false}, nil
}

func isSuccessful(code int) bool {
	return code >= http.StatusOK && code < http.StatusBadRequest
}
