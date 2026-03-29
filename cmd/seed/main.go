package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/juevigrace/diva-server/internal/di"
	"github.com/juevigrace/diva-server/internal/mail"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/server"
	"github.com/juevigrace/diva-server/storage"

	"github.com/joho/godotenv"
)

func main() {
	dev := flag.Bool("d", false, "Load .env.dev instead of .env")
	flag.Parse()

	envFile := ".env"
	if *dev {
		envFile = ".env.dev"
	}

	if err := godotenv.Load(envFile); err != nil {
		log.Fatal(err)
	}

	if err := os.Setenv(storage.DB_HOST_KEY, "localhost"); err != nil {
		log.Fatal(err)
	}

	if err := os.Setenv(storage.DB_PORT_KEY, "5433"); err != nil {
		log.Fatal(err)
	}

	conf := storage.NewDatabaseConf()
	database, err := storage.New(conf)
	if err != nil {
		log.Fatal("failed to create storage: %w", err)
	}

	queries := database.Queries()

	serverConf, ok := server.NewServerConfig().(*server.ServerConfig)
	if !ok {
		log.Fatal("invalid config")
	}

	mailClient := mail.NewClient(serverConf.ResendAPIKey, serverConf.ResendFromEmail)

	repoModule := di.NewRepoModule(queries)
	serviceModule := di.NewServiceModule(repoModule, mailClient)

	user := dtos.CreateUserDto{
		Email:    serverConf.RootEmail,
		Username: serverConf.RootUsername,
		Password: serverConf.RootPassword,
		Alias:    serverConf.RootUsername,
	}

	id, err := serviceModule.User.Create(context.Background(), &user)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(id)
}
