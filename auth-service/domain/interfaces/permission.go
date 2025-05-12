package interfaces

import (
	"auth-service/domain/entities"
	"context"
)

type PermissionUseCase interface {
	GetRoleAndPermissionByRoleUUID(ctx context.Context, roleUUID string) (*entities.Role, []entities.Permission, error)
}

type PermissionRepository interface {
	GetPermissionByUUID(ctx context.Context, uuid string) (*entities.Permission, error)
}

type RoleRepository interface {
	GetRoleByUUID(ctx context.Context, uuid string) (*entities.Role, error)
}

type RolePermissionRepository interface {
	GetRolePermissionByRoleUUID(roleUUID string) ([]entities.RolePermission, error)
	GetRolePermissionByPermissionUUID(permissionUUID string) ([]entities.RolePermission, error)
}
