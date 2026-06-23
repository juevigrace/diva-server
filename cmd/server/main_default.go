//go:build !postgres

package main

import (
	"context"
	"log"

	_ "github.com/joho/godotenv/autoload"
	"github.com/juevigrace/diva-server/server"
	"github.com/juevigrace/diva-server/storage/sqlite"
)

func main() {
	cfg := server.NewServerConfig()
	dbCfg := sqlite.NewSQLiteConf()
	dbCfg.LoadFromEnv()

	database, err := sqlite.New(dbCfg.(*sqlite.SQLiteConf))
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
