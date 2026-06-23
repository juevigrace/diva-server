package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/juevigrace/diva-server/pkg/concurrency"
	"github.com/juevigrace/diva-server/storage"
	pg "github.com/juevigrace/diva-server/storage/postgres/db"
	"github.com/juevigrace/diva-server/storage/postgres/db/migrations"
)

type PGStorage struct {
	pool    *pgxpool.Pool
	queries *pg.Queries
	config  *PGConf
}

func New(cfg *PGConf) (storage.Storage[pg.Queries], error) {
	dbInstance := new(PGStorage)

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("storage: %w", err)
	}

	dbInstance.config = cfg
	log.Println("Finished configuring")

	if err := dbInstance.initialize(); err != nil {
		return nil, err
	}

	return dbInstance, nil
}

func (s *PGStorage) initialize() error {
	if err := s.migrate(context.Background()); err != nil {
		return fmt.Errorf("storage: %v", err)
	}

	if err := s.openConnection(context.Background()); err != nil {
		return fmt.Errorf("storage: %v", err)
	}

	s.queries = pg.New(s.pool)

	return nil
}

func (s *PGStorage) migrate(ctx context.Context) error {
	return concurrency.WithTimeout(ctx, 5*time.Minute, func(ctx context.Context) error {
		url, err := s.config.Url()
		if err != nil {
			return fmt.Errorf("url: %v", err)
		}

		db, err := sql.Open("pgx", url)
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

func (s *PGStorage) openConnection(ctx context.Context) error {
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

func (s *PGStorage) Queries() *pg.Queries {
	return s.queries
}

func (s *PGStorage) Close() error {
	s.pool.Close()
	return nil
}
