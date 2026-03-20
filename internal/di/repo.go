package di

import (
	"github.com/juevigrace/diva-server/internal/repo"
	"github.com/juevigrace/diva-server/storage/db"
)

type RepoModule struct {
	User            *repo.UserRepository
	Session         *repo.SessionRepository
	UserPreferences *repo.UserPreferencesRepository
	UserPermission  *repo.UserPermissionRepository
	Verification    *repo.VerificationRepository
}

func NewRepoModule(queries *db.Queries) *RepoModule {
	return &RepoModule{
		User:            repo.NewUserRepository(queries),
		Session:         repo.NewSessionRepository(queries),
		UserPreferences: repo.NewUserPreferencesRepository(queries),
		UserPermission:  repo.NewUserPermissionRepository(queries),
		Verification:    repo.NewVerificationRepository(queries),
	}
}
