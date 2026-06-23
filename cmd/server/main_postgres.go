//go:build postgres

package main

import (
	"context"
	"log"

	_ "github.com/joho/godotenv/autoload"
	"github.com/juevigrace/diva-server/server"
	"github.com/juevigrace/diva-server/storage/postgres"
)

func main() {
	cfg := server.NewServerConfig()
	dbCfg := postgres.NewPGConf()
	dbCfg.LoadFromEnv()

	database, err := postgres.New(dbCfg.(*postgres.PGConf))
	if err != nil {
		log.Fatal(err)
	}

	newServer, err := server.NewServer(cfg, database)
	if err != nil {
		log.Fatal(err)
	}

	if err := newServer.ListenAndServe(context.Background()); err != nil {
		log.Fatal(err)
	}
}
