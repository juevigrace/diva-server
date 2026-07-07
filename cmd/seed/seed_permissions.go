package main

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/storage"
)

type seedPermission struct {
	Action    models.PermissionAction
	Name      string
	Desc      string
	RoleLevel models.Role
}

var seedPermissions = []seedPermission{
	{models.PERMISSION_USERS_READ, "Users Read", "View user accounts", models.ROLE_MODERATOR},
	{models.PERMISSION_USERS_WRITE, "Users Write", "Create and edit user accounts", models.ROLE_MODERATOR},
	{models.PERMISSION_USERS_EMAIL_WRITE, "Users Email Update", "Update user email addresses", models.ROLE_MODERATOR},
	{models.PERMISSION_USERS_PHONE_WRITE, "Users Phone Update", "Update user phone numbers", models.ROLE_MODERATOR},
	{models.PERMISSION_USERS_USERNAME_WRITE, "Users Username Update", "Update usernames", models.ROLE_MODERATOR},
	{models.PERMISSION_USERS_PASSWORD_WRITE, "Users Password Update", "Update user passwords", models.ROLE_MODERATOR},
	{models.PERMISSION_USERS_ROLE_WRITE, "Users Role Update", "Change user roles", models.ROLE_ADMIN},
	{models.PERMISSION_USERS_VERIFIED_WRITE, "Users Verified Update", "Verify user accounts", models.ROLE_ADMIN},
	{models.PERMISSION_USERS_RESTORE_WRITE, "Users Restore Write", "Restore deleted user accounts", models.ROLE_ADMIN},
	{models.PERMISSION_USERS_PROFILE_READ, "Users Profile Read", "View user profiles", models.ROLE_USER},
	{models.PERMISSION_USERS_PROFILE_WRITE, "Users Profile Write", "Create and update user profiles", models.ROLE_USER},
	{models.PERMISSION_USERS_PREFERENCES_READ, "Users Preferences Read", "View user preferences", models.ROLE_USER},
	{models.PERMISSION_USERS_PREFERENCES_WRITE, "Users Preferences Write", "Create and update user preferences", models.ROLE_USER},
	{models.PERMISSION_SESSIONS_READ, "Sessions Read", "View user sessions", models.ROLE_MODERATOR},
	{models.PERMISSION_SESSIONS_WRITE, "Sessions Write", "Manage user sessions", models.ROLE_MODERATOR},
	{models.PERMISSION_ACTIONS_READ, "Actions Read", "View user actions", models.ROLE_MODERATOR},
	{models.PERMISSION_ACTIONS_WRITE, "Actions Write", "Manage user actions", models.ROLE_MODERATOR},
	{models.PERMISSION_USER_PERMISSIONS_READ, "User Permissions Read", "View user permission grants", models.ROLE_ADMIN},
	{models.PERMISSION_USER_PERMISSIONS_WRITE, "User Permissions Write", "Grant and revoke user permissions", models.ROLE_ADMIN},
	{models.PERMISSION_PERMISSIONS_READ, "Permissions Read", "View all available permissions", models.ROLE_ADMIN},
	{models.PERMISSION_PERMISSIONS_WRITE, "Permissions Write", "Create and edit permissions", models.ROLE_ADMIN},
	{models.PERMISSION_OWNERSHIP_BYPASS, "Ownership Bypass", "Bypass resource ownership checks", models.ROLE_ADMIN},
}

func seedAllPermissions(ctx context.Context, store storage.PermissionStore) {
	for _, p := range seedPermissions {
		params := &storage.CreatePermissionParams{
			ID:          uuid.New(),
			Name:        p.Name,
			Description: p.Desc,
			Action:      p.Action.String(),
			RoleLevel:   p.RoleLevel.ToDB(),
		}
		if err := store.CreatePermission(ctx, params); err != nil {
			log.Fatalf("failed to seed permission %s: %v", p.Action.String(), err)
		}
	}

	log.Println("seeded all permissions")
}
