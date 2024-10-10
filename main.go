package main

import (
	"log"

	"github.com/juandr89/delivery-notifier-buyer/app_init"
	"github.com/juandr89/delivery-notifier-buyer/server"
)

func main() {
	cfg, err := server.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	app_init.RunServer(cfg)
}
