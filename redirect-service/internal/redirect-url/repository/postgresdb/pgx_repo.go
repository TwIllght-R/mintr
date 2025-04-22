package postgresdb

import (
	"context"
	"errors"
	"log"
	"redirect-service/domain/entities"
	"redirect-service/domain/interfaces"

	"gorm.io/gorm"
)

type redirectURLRepo struct {
	db *gorm.DB
}

func NewRedirectURLRepository(db *gorm.DB) interfaces.RedirectURLRepository {
	return &redirectURLRepo{db: db}
}

func (r *redirectURLRepo) GetByShortCode(ctx context.Context, shortCode string) (*entities.RedirectURL, error) {
	var redirectURL entities.RedirectURL
	err := r.db.Table("url_mappings").Where("short_code = ?", shortCode).First(&redirectURL).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entities.ErrURLNotFound
		}
		log.Println("Error fetching redirect URL :", err)
		return nil, entities.ErrInternalServer
	}
	return &redirectURL, nil
}
