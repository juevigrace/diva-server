//go:build !postgres

package main

import (
	"context"
	"flag"
	"log"

	"github.com/juevigrace/diva-server/internal/api/core/permission"
	"github.com/juevigrace/diva-server/internal/api/core/session"
	"github.com/juevigrace/diva-server/internal/api/core/user"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/server"
	"github.com/juevigrace/diva-server/storage/sqlite"

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

	dbCfg := sqlite.NewSQLiteConf()
	dbCfg.LoadFromEnv()

	database, err := sqlite.New(dbCfg.(*sqlite.SQLiteConf))
	if err != nil {
		log.Fatal("failed to create storage: %w", err)
	}

	serverConf := server.NewServerConfig().(*server.ServerConfig)

	pModule := permission.NewPermissionModule(database.PermissionStore())
	sModule := session.NewSessionModule(database.SessionStore())
	uModule := user.NewUserModule(
		database.UserStore(),
		database.UserActionStore(),
		database.UserPermissionStore(),
		database.UserPreferenceStore(),
		database.UserProfileStore(),
		database.UserStateStore(),
		pModule.Repo,
		sModule.Repo,
		sModule.Handler,
		nil,
	)

	userDto := dtos.CreateUserDto{
		Email:    serverConf.RootEmail,
		Username: serverConf.RootUsername,
		Password: serverConf.RootPassword,
	}

	id, err := uModule.URepo.Create(context.Background(), &userDto)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(id)
}
