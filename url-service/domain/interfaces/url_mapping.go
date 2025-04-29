package interfaces

import (
	"context"
	"url-service/domain/entities"
	"url-service/domain/events"
	"url-service/domain/pagination"
)

type URLMappingUseCase interface {
	GenerateShortCode(ctx context.Context, ownerUUID string, in entities.URLMapping) (*entities.URLMapping, error)
	GetURLMappingByUUID(ctx context.Context, uuid, ownerUUID string) (*entities.URLMapping, error)
	GetAllURLMapping(ctx context.Context, ownerUUID string, in pagination.Pagination) (*pagination.PaginatedResponse[entities.URLMapping], error)
	DeleteURLMapping(ctx context.Context, uuid, ownerUUID string) error
	UpdateURLMapping(ctx context.Context, uuid, ownerUUID string, in entities.URLMapping) error
}

type URLMappingRepository interface {
	Store(ctx context.Context, in entities.URLMapping) error
	GetByUUID(ctx context.Context, uuid, ownerUUID string) (*entities.URLMapping, error)
	Delete(ctx context.Context, uuid, ownerUUID string) error
	Update(ctx context.Context, uuid, ownerUUID string, in entities.URLMapping) error
	GetAll(ctx context.Context, ownerUUID string, in pagination.Pagination) ([]entities.URLMapping, int64, error)
	Fetch(ctx context.Context, ownerUUID string, in pagination.CursorPagination) ([]entities.URLMapping, *string, error)
}

type URLMappingEventProducer interface {
	ProduceURLMappingCreatedEvent(ctx context.Context, in events.URLCreatedEvent) error
	ProduceURLMappingDeletedEvent(ctx context.Context, in events.URLDeletedEvent) error
	ProduceURLMappingUpdatedEvent(ctx context.Context, in events.URLUpdatedEvent) error
}
