package web

import (
	"btcRate/docs"
	"btcRate/infrastructure"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
)

type IExtendedHttpClient interface {
	SendRequest(req *http.Request) ([]byte, int, error)
}

type ServerManager struct {
	client IExtendedHttpClient
}

type Response[T any] struct {
	Code         int
	Body         *T
	ErrorMessage string
	Successful   bool
}

func NewServerManager() ServerManager {
	return ServerManager{infrastructure.NewExtendedHttpClient(nil)}
}

func (*ServerManager) RunServer(storageFile string) (func() error, error) {
	r := gin.Default()
	r.Use(errorHandlingMiddleware())

	btcUahController, err := newBtcUahController(storageFile)
	if err != nil {
		return nil, err
	}

	docs.SwaggerInfo.BasePath = apiBasePath
	api := r.Group(apiBasePath)
	{
		api.GET(getRate, btcUahController.getRate)
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

	body, statusCode, err := s.client.SendRequest(req)
	if err != nil {
		return nil, err
	}

	if isSuccessful(statusCode) {
		var result int
		err = json.Unmarshal(body, &result)
		if err != nil {
			return nil, err
		}
		return &Response[int]{Code: statusCode, Body: &result, ErrorMessage: "", Successful: true}, nil
	}

	return &Response[int]{Code: statusCode, ErrorMessage: string(body), Successful: false}, nil
}

func isSuccessful(code int) bool {
	return code >= http.StatusOK && code < http.StatusBadRequest
}
