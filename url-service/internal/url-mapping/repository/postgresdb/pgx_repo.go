package postgresdb

import (
	"context"
	"errors"
	"fmt"
	"log"
	"url-service/domain/entities"
	"url-service/domain/interfaces"
	"url-service/domain/pagination"

	"gorm.io/gorm"
)

type urlMappingRepo struct {
	db *gorm.DB
}

func NewURLMappingRepo(db *gorm.DB) interfaces.URLMappingRepository {
	return &urlMappingRepo{db: db}
}

func (r *urlMappingRepo) Store(ctx context.Context, in entities.URLMapping) error {
	err := r.db.Create(&in).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return entities.ErrURLAlreadyExists
		}
		log.Println("Error storing URL mapping:", err)
		return entities.ErrInternalServer
	}
	return nil
}

func (r *urlMappingRepo) GetByUUID(ctx context.Context, uuid string) (*entities.URLMapping, error) {
	var urlMapping entities.URLMapping
	err := r.db.Where("uuid = ?", uuid).First(&urlMapping).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entities.ErrURLNotFound
		}
		log.Println("Error fetching URL mapping:", err)
		return nil, entities.ErrInternalServer
	}
	return &urlMapping, nil
}

func (r *urlMappingRepo) Delete(ctx context.Context, uuid string) error {
	result := r.db.Where("uuid = ?", uuid).Delete(&entities.URLMapping{})
	if result.Error != nil {
		log.Println("Error deleting URL mapping:", result.Error)
		return entities.ErrInternalServer
	}
	if result.RowsAffected == 0 {
		return entities.ErrURLNotFound
	}
	return nil
}

func (r *urlMappingRepo) Update(ctx context.Context, uuid string, in entities.URLMapping) error {
	result := r.db.Model(&entities.URLMapping{}).Where("uuid = ?", uuid).Updates(in)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return entities.ErrURLAlreadyExists
		}
		log.Println("Error updating URL mapping:", result.Error)
		return entities.ErrInternalServer
	}
	if result.RowsAffected == 0 {
		return entities.ErrURLNotFound
	}
	return nil
}

func (r *urlMappingRepo) GetAll(ctx context.Context, in pagination.Pagination) ([]entities.URLMapping, int64, error) {
	var urlMappings []entities.URLMapping
	var total int64

	err := r.db.Model(&entities.URLMapping{}).Count(&total).Error
	if err != nil {
		log.Println("Error counting URL mappings:", err)
		return nil, 0, entities.ErrInternalServer
	}

	err = r.db.Order(fmt.Sprintf("%s %s", in.SortBy, in.Sort)).Offset((in.Page - 1) * in.PageSize).Limit(in.PageSize).Find(&urlMappings).Error
	if err != nil {
		log.Println("Error fetching URL mappings:", err)
		return nil, 0, entities.ErrInternalServer
	}
	return urlMappings, total, nil
}

func (r *urlMappingRepo) Fetch(ctx context.Context, in pagination.CursorPagination) ([]entities.URLMapping, *string, error) {
	var urlMappings []entities.URLMapping
	var nextCursor *string

	query := r.db.Model(&entities.URLMapping{}).Limit(in.Limit)

	if in.Cursor != nil {
		query = query.Where("uuid > ?", *in.Cursor)
	}

	err := query.Find(&urlMappings).Error
	if err != nil {
		log.Println("Error fetching URL mappings:", err)
		return nil, nil, entities.ErrInternalServer
	}

	if len(urlMappings) > 0 {
		nextCursor = &urlMappings[len(urlMappings)-1].UUID
	}

	return urlMappings, nextCursor, nil
}
