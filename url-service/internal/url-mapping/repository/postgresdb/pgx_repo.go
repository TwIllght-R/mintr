package postgresdb

import (
	"context"
	"errors"
	"fmt"
	"url-service/domain/entities"
	"url-service/domain/interfaces"
	"url-service/domain/pagination"

	"gorm.io/gorm"
)

type urlMappingRepository struct {
	db *gorm.DB
}

func NewURLMappingRepository(db *gorm.DB) interfaces.URLMappingRepository {
	return &urlMappingRepository{db: db}
}

func (r *urlMappingRepository) Store(ctx context.Context, in entities.URLMapping) error {
	err := r.db.Table("url_mappings").Create(&in).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return entities.ErrURLAlreadyExists
		}
		return fmt.Errorf("error storing URL mapping: %w", err)
	}
	return nil
}

func (r *urlMappingRepository) GetByUUID(ctx context.Context, uuid, ownerUUID string) (*entities.URLMapping, error) {
	var urlMapping entities.URLMapping
	err := r.db.Table("url_mappings").Where("uuid = ? AND owner_uuid = ?", uuid, ownerUUID).First(&urlMapping).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entities.ErrURLNotFound
		}

		return nil, fmt.Errorf("error fetching URL mapping: %w", err)
	}
	return &urlMapping, nil
}

func (r *urlMappingRepository) Delete(ctx context.Context, uuid, ownerUUID string) error {
	result := r.db.Table("url_mappings").Where("uuid = ? AND owner_uuid = ?", uuid, ownerUUID).Delete(&entities.URLMapping{})
	if result.Error != nil {
		return fmt.Errorf("error deleting URL mapping: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return entities.ErrURLNotFound
	}
	return nil
}

func (r *urlMappingRepository) Update(ctx context.Context, uuid, ownerUUID string, in entities.URLMapping) error {
	result := r.db.Table("url_mappings").Where("uuid = ? AND owner_uuid = ?", uuid, ownerUUID).Updates(in)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return entities.ErrURLAlreadyExists
		}
		return fmt.Errorf("error updating URL mapping: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return entities.ErrURLNotFound
	}
	return nil
}

func (r *urlMappingRepository) GetAll(ctx context.Context, ownerUUID string, in pagination.Pagination) ([]entities.URLMapping, int64, error) {
	var urlMappings []entities.URLMapping
	var total int64

	err := r.db.Table("url_mappings").Where("owner_uuid = ?", ownerUUID).Model(&entities.URLMapping{}).Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("error counting URL mappings: %w", err)
	}

	err = r.db.Table("url_mappings").Where("owner_uuid = ?", ownerUUID).Order(fmt.Sprintf("%s %s", in.SortBy, in.Sort)).Offset((in.Page - 1) * in.PageSize).Limit(in.PageSize).Find(&urlMappings).Error
	if err != nil {
		return nil, 0, fmt.Errorf("error fetching URL mappings: %w", err)
	}
	return urlMappings, total, nil
}

func (r *urlMappingRepository) Fetch(ctx context.Context, ownerUUID string, in pagination.CursorPagination) ([]entities.URLMapping, *string, error) {
	var urlMappings []entities.URLMapping
	var nextCursor *string

	query := r.db.Table("url_mappings").Where("owner_uuid = ?", ownerUUID).Model(&entities.URLMapping{}).Limit(in.Limit)

	if in.Cursor != nil {
		query = query.Where("uuid > ?", *in.Cursor)
	}

	err := query.Find(&urlMappings).Error
	if err != nil {
		return nil, nil, fmt.Errorf("error fetching URL mappings: %w", err)
	}

	if len(urlMappings) > 0 {
		nextCursor = &urlMappings[len(urlMappings)-1].UUID
	}

	return urlMappings, nextCursor, nil
}
