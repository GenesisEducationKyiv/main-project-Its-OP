package main

import (
	"btcRate/coin/infrastructure/bus"
	"btcRate/coin/web"
	"context"
	"log"
	"os"
	"os/signal"
)

func main() {
	busConfig := &bus.RabbitMQConfig{
		Host:     os.Getenv("RABBITMQ_HOST"),
		User:     os.Getenv("RABBITMQ_USER"),
		Password: os.Getenv("RABBITMQ_PASSWORD"),
	}
	commandBus, router := bus.AddCommandBus(busConfig)
	go func() {
		_ = router.Run(context.Background())
	}()

	server := web.NewServerManager(commandBus)
	stop, err := server.RunServer("./logs/coin-logs.csv")
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	if err := stop(); err != nil {
		log.Fatal("Failed to stop the server: ", err)
	}

	if err := router.Close(); err != nil {
		log.Fatal("Failed to stop the Router: ", err)
	}
}
