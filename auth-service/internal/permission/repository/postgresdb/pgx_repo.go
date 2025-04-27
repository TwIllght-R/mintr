package postgresdb

import (
	"auth-service/domain/entities"
	"auth-service/domain/interfaces"
	"context"
	"errors"
	"log"

	"gorm.io/gorm"
)

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) interfaces.PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) GetPermissionByUUID(ctx context.Context, uuid string) (*entities.Permission, error) {
	var permission entities.Permission
	err := r.db.Table("permissions").Where("uuid = ?", uuid).First(&permission).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entities.ErrPermissionNotFound
		}
		log.Println("Error fetching permission:", err)
		return nil, entities.ErrInternalServer
	}
	return &permission, nil
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) interfaces.RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) GetRoleByUUID(ctx context.Context, uuid string) (*entities.Role, error) {
	var role entities.Role
	err := r.db.Table("roles").Where("uuid = ?", uuid).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entities.ErrRoleNotFound
		}
		log.Println("Error fetching role:", err)
		return nil, entities.ErrInternalServer
	}
	return &role, nil
}

type rolePermissionRepository struct {
	db *gorm.DB
}

func NewRolePermissionRepository(db *gorm.DB) interfaces.RolePermissionRepository {
	return &rolePermissionRepository{db: db}
}

func (r *rolePermissionRepository) GetRolePermissionByRoleUUID(roleUUID string) ([]entities.RolePermission, error) {
	var rolePermissions []entities.RolePermission
	err := r.db.Table("role_permissions").Where("role_uuid = ?", roleUUID).Find(&rolePermissions).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entities.ErrRolePermissionNotFound
		}
		log.Println("Error fetching role permissions by role UUID:", err)
		return nil, entities.ErrInternalServer
	}
	return rolePermissions, nil
}

func (r *rolePermissionRepository) GetRolePermissionByPermissionUUID(permissionUUID string) ([]entities.RolePermission, error) {
	var rolePermissions []entities.RolePermission
	err := r.db.Table("role_permissions").Where("permission_uuid = ?", permissionUUID).Find(&rolePermissions).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entities.ErrRolePermissionNotFound
		}
		log.Println("Error fetching role permissions by permission UUID:", err)
		return nil, entities.ErrInternalServer
	}
	return rolePermissions, nil
}
