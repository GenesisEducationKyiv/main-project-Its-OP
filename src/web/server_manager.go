package web

import (
	"btcRate/docs"
	"context"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
)

type ServerManager struct {
}

func (*ServerManager) RunServer(storageFile string) (func() error, error) {
	r := gin.Default()
	r.Use(errorHandlingMiddleware())

	btcUahController, err := newBtcUahController(storageFile)
	if err != nil {
		return nil, err
	}

	docs.SwaggerInfo.BasePath = "/api/v1"
	api := r.Group("/api/v1")
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

	return stop, err
}
