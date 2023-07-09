package main

import (
	"btcRate/campaign/web"
	"log"
	"os"
	"os/signal"
)

func main() {
	server := web.NewServerManager()

	fc := &web.FileConfiguration{LogStorageFile: "./data/logs.csv", EmailStorageFile: "./data/emails.json"}
	sc := &web.SendgridConfiguration{ApiKey: os.Getenv("SENDGRID_KEY"), SenderName: os.Getenv("SENDGRID_SENDER_NAME"), SenderEmail: os.Getenv("SENDGRID_SENDER_EMAIL")}
	pc := &web.ProviderConfiguration{Hostname: os.Getenv("COIN_HOST"), Schema: os.Getenv("COIN_SCHEMA")}

	stop, err := server.RunServer(fc, sc, pc)
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
