package usecase

import (
	"auth-service/domain/entities"
	"auth-service/domain/interfaces"
	"context"
)

type permissionUseCase struct {
	permissionRepo     interfaces.PermissionRepository
	roleRepo           interfaces.RoleRepository
	rolePermissionRepo interfaces.RolePermissionRepository
}

func NewPermissionUsecase(
	permissionRepo interfaces.PermissionRepository,
	roleRepo interfaces.RoleRepository,
	rolePermissionRepo interfaces.RolePermissionRepository,
) interfaces.PermissionUseCase {
	return &permissionUseCase{
		permissionRepo:     permissionRepo,
		roleRepo:           roleRepo,
		rolePermissionRepo: rolePermissionRepo,
	}
}

func (u *permissionUseCase) GetRoleAndPermissionByRoleUUID(ctx context.Context, roleUUID string) (*entities.Role, []entities.Permission, error) {
	role, err := u.roleRepo.GetRoleByUUID(ctx, roleUUID)
	if err != nil {
		return nil, nil, err
	}

	rolePermissions, err := u.rolePermissionRepo.GetRolePermissionByRoleUUID(roleUUID)
	if err != nil {
		return nil, nil, err
	}

	var permissions []entities.Permission
	for _, rolePermission := range rolePermissions {
		permission, err := u.permissionRepo.GetPermissionByUUID(ctx, rolePermission.PermissionUUID)
		if err != nil {
			return nil, nil, err
		}
		permissions = append(permissions, *permission)
	}

	return role, permissions, nil
}
