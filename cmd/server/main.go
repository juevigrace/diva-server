package main

import (
	"context"
	"log"

	_ "github.com/joho/godotenv/autoload"
	"github.com/juevigrace/diva-server/server"
)

func main() {
	config := server.NewServerConfig()

	newServer, err := server.NewServer(config)
	if err != nil {
		log.Fatal(err)
	}

	if err := newServer.ListenAndServe(context.Background()); err != nil {
		log.Fatal(err)
	}
}
