package main

import (
	"btcRate/web"
	"log"
)

func main() {
	err := web.RunBtcUahController()
	if err == nil {
		log.Fatal(err)
	}
}
