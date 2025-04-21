package interfaces

import (
	"url-service/domain/entities"
	"url-service/domain/pagination"
	"context"
)

type URLMappingUseCase interface {
	GenerateShortCode(ctx context.Context, in entities.URLMapping) (*entities.URLMapping, error)
	GetURLMappingByUUID(ctx context.Context, uuid string) (*entities.URLMapping, error)
	GetAllURLMapping(ctx context.Context, in pagination.Pagination) (*pagination.PaginatedResponse[entities.URLMapping], error)
	DeleteURLMapping(ctx context.Context, uuid string) error
	UpdateURLMapping(ctx context.Context, uuid string, in entities.URLMapping) error
}

type URLMappingRepository interface {
	Store(ctx context.Context, in entities.URLMapping) error
	GetByUUID(ctx context.Context, uuid string) (*entities.URLMapping, error)
	Delete(ctx context.Context, uuid string) error
	Update(ctx context.Context, uuid string, in entities.URLMapping) error
	GetAll(ctx context.Context, in pagination.Pagination) ([]entities.URLMapping, int64, error)
	Fetch(ctx context.Context, in pagination.CursorPagination) ([]entities.URLMapping, *string, error)
}
