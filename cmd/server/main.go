package main

import (
	"go-payments-api/di"
	"log"
)

// @title Microservice Payments API
// @version 1.0
// @description This is a sample server for a microservice payments API.

// @Contact.url http://www.example.com/support
// @Contact.name API Support

// @host localhost:8080
// @BasePath /v1/payments
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-KEY
func main() {
	// sdk, err := distro.Run()
	// if err != nil {
	// 	log.Fatalf("Failed to run distro: %v", err)
	// }
	// defer func() {
	// 	if err := sdk.Shutdown(context.Background()); err != nil {
	// 		log.Fatalf("Failed to shutdown SDK: %v", err)
	// 	}
	// }()

	api, cleanup, err := di.InitializeApi()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	api.Start()
	defer cleanup()
}
