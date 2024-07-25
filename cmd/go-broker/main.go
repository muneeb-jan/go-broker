package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/muneeb-jan/go-broker/internal/controller"
	"github.com/muneeb-jan/go-broker/internal/messagebroker"
)

func main() {
	devMode := flag.Bool("dev", false, "Run in development mode")
	flag.Parse()

	broker := messagebroker.NewBroker()
	ctrl := controller.NewController(broker, *devMode)

	server := &http.Server{
		Addr:    ":8080",
		Handler: ctrl.Routes(),
	}

	if *devMode {
		log.Printf("Server running on port 8080 (dev mode: %v)\n", *devMode)
	} else {
		log.Println("Server running on port 8080")
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
