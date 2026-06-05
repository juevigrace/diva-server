package middlewares

import "github.com/juevigrace/diva-server/internal/models"

type RoleRequirement struct {
	Satisfied bool
	Roles     []models.Role
}

type PermissionRequirement struct {
	Satisfied         bool
	PermissionActions []models.PermissionAction
}

type OwnershipRequirement struct {
	Satisfied bool
	Cache     map[string]any
}
