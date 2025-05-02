package usecase

import (
	"analytics-service/domain/entities"
	"analytics-service/domain/interfaces"
	"context"
	"log"
)

type analyticsUseCase struct {
	analyticsRepository interfaces.AnalyticsRepository
}

func NewAnalyticsUseCase(analyticsRepository interfaces.AnalyticsRepository) interfaces.AnalyticsUseCase {
	return &analyticsUseCase{
		analyticsRepository: analyticsRepository,
	}
}

func (u *analyticsUseCase) RecordClicked(ctx context.Context, in entities.URLClicked) error {
	err := u.analyticsRepository.Store(ctx, in)
	if err != nil {
		log.Println(err)
		return entities.ErrInternalServer
	}
	return nil
}

func (u *analyticsUseCase) GetClickedDetails(ctx context.Context, urlUUID, ownerUUID string) ([]entities.URLClicked, error) {
	urls, err := u.analyticsRepository.GetAllByUrlUUID(ctx, urlUUID, ownerUUID)
	if err != nil {
		log.Println(err)
		return nil, entities.ErrInternalServer
	}
	return urls, nil
}
