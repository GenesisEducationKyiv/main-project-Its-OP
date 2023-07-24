package main

import (
	"btcRate/campaign/infrastructure/bus"
	"btcRate/campaign/web"
	"context"
	"log"
	"os"
	"os/signal"
)

func main() {
	commandBus, router, err := bus.AddCommandBus(os.Getenv("KAFKA_HOST"), "campaign-consumer-group")
	if err != nil {
		log.Fatal(err.Error())
	}

	go func() {
		if err := router.Run(context.Background()); err != nil {
			log.Printf("Error running router: %v", err)
		}
	}()

	server := web.NewServerManager()

	fc := &web.FileConfiguration{EmailStorageFile: "./data/emails.json"}
	sc := &web.SendgridConfiguration{ApiKey: os.Getenv("SENDGRID_KEY"), SenderName: os.Getenv("SENDGRID_SENDER_NAME"), SenderEmail: os.Getenv("SENDGRID_SENDER_EMAIL")}
	pc := &web.ProviderConfiguration{Hostname: os.Getenv("COIN_HOST"), Schema: os.Getenv("COIN_SCHEMA")}

	stop, err := server.RunServer(fc, sc, pc, commandBus)
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	if err := stop(); err != nil {
		log.Fatal("Failed to stop the server: ", err)
	}
}
