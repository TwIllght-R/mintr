package postgresdb

import (
	"auth-service/domain/entities"
	"auth-service/domain/interfaces"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type customerAuthRepository struct {
	db *gorm.DB
}

func NewCustomerAuthRepository(db *gorm.DB) interfaces.CustomerAuthRepository {
	return &customerAuthRepository{db: db}
}

func (r *customerAuthRepository) Store(ctx context.Context, in entities.CustomerAuth) error {
	err := r.db.Table("customer_auths").Create(&in).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return entities.ErrEmailAlreadyExists
		}
		return fmt.Errorf("failed to store customer auth: %w", err)
	}
	return nil
}

func (r *customerAuthRepository) Update(ctx context.Context, uuid string, in entities.CustomerAuth) error {
	err := r.db.Table("customer_auths").Where("uuid = ?", uuid).Updates(&in).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.ErrUserNotFound
		}
		return fmt.Errorf("failed to update customer auth: %w", err)
	}
	return nil
}

func (r *customerAuthRepository) GetByUUID(ctx context.Context, uuid string) (*entities.CustomerAuth, error) {
	var auth entities.CustomerAuth
	err := r.db.Table("customer_auths").Where("uuid = ?", uuid).First(&auth).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entities.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to fetch customer auth: %w", err)
	}
	return &auth, nil
}

func (r *customerAuthRepository) GetByEmail(ctx context.Context, email string) (*entities.CustomerAuth, error) {
	var auth entities.CustomerAuth
	err := r.db.Table("customer_auths").Where("email = ?", email).First(&auth).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entities.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to fetch customer auth: %w", err)
	}
	return &auth, nil
}
