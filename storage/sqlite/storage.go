package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/juevigrace/diva-server/pkg/concurrency"
	"github.com/juevigrace/diva-server/storage"
	sqli "github.com/juevigrace/diva-server/storage/sqlite/db"
	"github.com/juevigrace/diva-server/storage/sqlite/db/migrations"

	_ "modernc.org/sqlite"
)

type SQLiteStorage struct {
	db      *sql.DB
	queries *sqli.Queries
	config  *SQLiteConf
}

func New(cfg *SQLiteConf) (storage.Storage[sqli.Queries], error) {
	dbInstance := new(SQLiteStorage)

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

func (s *SQLiteStorage) initialize() error {
	if err := s.migrate(context.Background()); err != nil {
		return fmt.Errorf("storage: %v", err)
	}

	if err := s.openConnection(context.Background()); err != nil {
		return fmt.Errorf("storage: %v", err)
	}

	s.queries = sqli.New(s.db)

	return nil
}

func (s *SQLiteStorage) migrate(ctx context.Context) error {
	return concurrency.WithTimeout(ctx, 5*time.Minute, func(ctx context.Context) error {
		url, err := s.config.Url()
		if err != nil {
			return fmt.Errorf("url: %v", err)
		}

		db, err := sql.Open("sqlite", url)
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

func (s *SQLiteStorage) openConnection(ctx context.Context) error {
	url, err := s.config.Url()
	if err != nil {
		return fmt.Errorf("url: %v", err)
	}

	return concurrency.WithTimeout(ctx, 30*time.Second, func(ctx context.Context) error {
		db, err := sql.Open("sqlite", url)
		if err != nil {
			return fmt.Errorf("open: %v", err)
		}
		if err = db.PingContext(ctx); err != nil {
			return fmt.Errorf("ping: %v", err)
		}
		s.db = db
		return nil
	})
}

func (s *SQLiteStorage) Queries() *sqli.Queries {
	return s.queries
}

func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}
