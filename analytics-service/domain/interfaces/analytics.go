package interfaces

import (
	"analytics-service/domain/entities"
	"context"
)

type AnalyticsUseCase interface {
	RecordClicked(ctx context.Context, in entities.URLClicked) error
	GetClickedDetails(ctx context.Context, urlUUID, ownerUUID string) ([]entities.URLClicked, error)
}

type AnalyticsRepository interface {
	Store(ctx context.Context, in entities.URLClicked) error
	GetAllByUrlUUID(ctx context.Context, urlUUID, ownerUUID string) ([]entities.URLClicked, error)
}
