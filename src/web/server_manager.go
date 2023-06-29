package web

import (
	"btcRate/docs"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"io"
	"log"
	"net/http"
)

type ServerManager struct {
}

type Response[T any] struct {
	Code         int
	Body         *T
	ErrorMessage string
	Successful   bool
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
	body, statusCode, err := s.sendRequest(url)
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

func (*ServerManager) sendRequest(url string) ([]byte, int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, 0, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	if err = resp.Body.Close(); err != nil {
		return nil, 0, err
	}

	return body, resp.StatusCode, nil
}

func isSuccessful(code int) bool {
	return code >= http.StatusOK && code < http.StatusBadRequest
}
