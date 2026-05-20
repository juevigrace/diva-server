package di

import (
	"github.com/juevigrace/diva-server/internal/repo"
	"github.com/juevigrace/diva-server/storage/db"
)

type RepoModule struct {
	User            *repo.UserRepo
	Session         *repo.SessionRepo
	UserPreferences *repo.UserPreferencesRepo
	UserPermission  *repo.UserPermsRepo
	Verification    *repo.UserVerificationRepo
	Action          *repo.UserActionsRepo
	UserProfile     *repo.UserProfileRepo
	Permissions     *repo.PermissionsRepo
}

func NewRepoModule(queries *db.Queries) *RepoModule {
	return &RepoModule{
		User:            repo.NewUserRepo(queries),
		Session:         repo.NewSessionRepo(queries),
		UserPreferences: repo.NewUserPreferencesRepo(queries),
		UserPermission:  repo.NewUserPermsRepo(queries),
		Verification:    repo.NewUserVerificationRepo(queries),
		Action:          repo.NewUserActionsRepo(queries),
		UserProfile:     repo.NewUserProfileRepo(queries),
		Permissions:     repo.NewPermissionsRepo(queries),
	}
}
