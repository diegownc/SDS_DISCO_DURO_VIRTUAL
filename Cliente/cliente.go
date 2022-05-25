package main

import (
	"log"

	"cliente/api"
)

const (
	serverAddress = "localhost:8080"
)

func main() {
	server := api.NewServer()

	err := server.Start(serverAddress)
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}
}
