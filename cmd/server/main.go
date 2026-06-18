package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/juevigrace/diva-server/server"
)

func main() {
	flags := server.NewServerFlags()

	if flags.UsesEnv {
		for _, f := range []string{".env", ".env.dev"} {
			if _, err := os.Stat(f); err == nil {
				godotenv.Load(f)
			}
		}
		log.Printf("envs loaded\n")
	}

	config := server.NewServerConfig(flags)
	newServer, err := server.NewServer(config)
	if err != nil {
		log.Fatal(err)
	}

	if err := newServer.ListenAndServe(context.Background()); err != nil {
		log.Fatal(err)
	}
}
