package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/juevigrace/diva-server/concurrency"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/storage/db"
	"github.com/juevigrace/diva-server/storage/db/migrations"
)

type Storage interface {
	Close() error
	Health(ctx context.Context) HealthResult
	Queries() *db.Queries
}

type StorageS struct {
	pool    *pgxpool.Pool
	queries *db.Queries
	config  *DatabaseConf
}

func New(config models.Config) (Storage, error) {
	dbInstance := new(StorageS)
	dbInstance.config = NewDatabaseConf().(*DatabaseConf)
	dbInstance.config.Configure(config)
	log.Println("Finished configuring")

	if err := dbInstance.initialize(); err != nil {
		return nil, err
	}

	return dbInstance, nil
}

func (s *StorageS) initialize() error {
	if err := s.migrate(context.Background()); err != nil {
		return fmt.Errorf("storage: %v", err)
	}

	if err := s.openConnection(context.Background()); err != nil {
		return fmt.Errorf("storage: %v", err)
	}

	s.queries = db.New(s.pool)

	return nil
}

func (s *StorageS) migrate(ctx context.Context) error {
	return concurrency.WithTimeout(ctx, 5*time.Minute, func(ctx context.Context) error {
		url, err := s.config.Url()
		if err != nil {
			return fmt.Errorf("url: %v", err)
		}

		db, err := sql.Open(s.config.Driver, url)
		if err != nil {
			return fmt.Errorf("open connection: %v", err)
		}

		if err := migrations.RunMigrations(db); err != nil {
			return fmt.Errorf("migrate: %v", err)
		}

		if err := db.Close(); err != nil {
			return fmt.Errorf("close: %v", err)
		}

		return nil
	})
}

func (s *StorageS) openConnection(ctx context.Context) error {
	url, err := s.config.Url()
	if err != nil {
		return fmt.Errorf("url: %v", err)
	}

	return concurrency.WithTimeout(ctx, 30*time.Second, func(ctx context.Context) error {
		conn, err := pgxpool.New(ctx, url)
		if err != nil {
			return fmt.Errorf("open pool: %v", err)
		}
		if err = conn.Ping(ctx); err != nil {
			return fmt.Errorf("ping: %v", err)
		}
		s.pool = conn
		return nil
	})
}

func (s *StorageS) Queries() *db.Queries {
	return s.queries
}

func (s *StorageS) Close() error {
	s.pool.Close()
	return nil
}
