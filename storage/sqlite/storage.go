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
	db                 *sql.DB
	queries            *sqli.Queries
	config             *SQLiteConf
	userStore          *UserStore
	permissionStore    *PermissionStore
	sessionStore       *SessionStore
	userStateStore     *UserStateStore
	userProfileStore   *UserProfileStore
	userPreferenceStore *UserPreferenceStore
	userPermissionStore *UserPermissionStore
	userActionStore    *UserActionStore
	userVerificationStore *UserVerificationStore
}

func New(cfg *SQLiteConf) (storage.Storage, error) {
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

	s.userStore = NewUserStore(s.queries)
	s.permissionStore = NewPermissionStore(s.queries)
	s.sessionStore = NewSessionStore(s.queries)
	s.userStateStore = NewUserStateStore(s.queries)
	s.userProfileStore = NewUserProfileStore(s.queries)
	s.userPreferenceStore = NewUserPreferenceStore(s.queries)
	s.userPermissionStore = NewUserPermissionStore(s.queries)
	s.userActionStore = NewUserActionStore(s.queries)
	s.userVerificationStore = NewUserVerificationStore(s.queries)

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

func (s *SQLiteStorage) UserStore() storage.UserStore               { return s.userStore }
func (s *SQLiteStorage) PermissionStore() storage.PermissionStore    { return s.permissionStore }
func (s *SQLiteStorage) SessionStore() storage.SessionStore          { return s.sessionStore }
func (s *SQLiteStorage) UserStateStore() storage.UserStateStore      { return s.userStateStore }
func (s *SQLiteStorage) UserProfileStore() storage.UserProfileStore  { return s.userProfileStore }
func (s *SQLiteStorage) UserPreferenceStore() storage.UserPreferenceStore   { return s.userPreferenceStore }
func (s *SQLiteStorage) UserPermissionStore() storage.UserPermissionStore   { return s.userPermissionStore }
func (s *SQLiteStorage) UserActionStore() storage.UserActionStore    { return s.userActionStore }
func (s *SQLiteStorage) UserVerificationStore() storage.UserVerificationStore { return s.userVerificationStore }

func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}
