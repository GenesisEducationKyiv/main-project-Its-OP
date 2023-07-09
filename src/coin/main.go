package main

import (
	"btcRate/coin/web"
	"log"
	"os"
	"os/signal"
)

func main() {
	server := web.NewServerManager()
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
}
