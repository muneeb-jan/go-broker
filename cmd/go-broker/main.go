package main

import (
	"log"
	"net/http"

	"github.com/muneeb-jan/go-broker/internal/controller"
	"github.com/muneeb-jan/go-broker/internal/messagebroker"
)

func main() {
	broker := messagebroker.NewBroker()
	ctrl := controller.NewController(broker)

	server := &http.Server{
		Addr:    ":8080",
		Handler: ctrl.Routes(),
	}

	log.Println("Server running on port 8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
