package main

import (
	"log"

	"github.com/diegownc/SDS_DISCO_DURO_VIRTUAL/api"
)

const (
	serverAddress = "localhost:8081"
)

func main() {
	server := api.NewServer()

	err := server.Start(serverAddress)
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}
}
