package usecase

import (
	"context"
	"crypto/rand"
	"errors"
	"log"
	"math/big"
	"time"
	"url-service/domain/entities"
	"url-service/domain/events"
	"url-service/domain/interfaces"
	"url-service/domain/pagination"

	"github.com/google/uuid"
)

type URLMappingUseCase struct {
	urlMappingRepo     interfaces.URLMappingRepository
	urlMappingProducer interfaces.URLMappingEventProducer
}

func NewURLMappingUseCase(urlMappingRepo interfaces.URLMappingRepository, urlMappingProducer interfaces.URLMappingEventProducer) interfaces.URLMappingUseCase {
	return &URLMappingUseCase{
		urlMappingRepo:     urlMappingRepo,
		urlMappingProducer: urlMappingProducer,
	}
}

func (u *URLMappingUseCase) GenerateShortCode(ctx context.Context, ownerUUID string, in entities.URLMapping) (*entities.URLMapping, error) {
	// Generate a unique short code
	shortCode, err := generateShortCode()
	if err != nil {
		log.Println("Error generating short code:", err)
		return nil, entities.ErrInternalServer
	}

	urlMapping := entities.URLMapping{
		UUID:        uuid.New().String(),
		OwnerUUID:   ownerUUID,
		ShortCode:   shortCode,
		OriginalURL: in.OriginalURL,
		Title:       in.Title,
		UTMSource:   in.UTMSource,
		UTMMedium:   in.UTMMedium,
		UTMCampaign: in.UTMCampaign,
		UTMTerm:     in.UTMTerm,
		UTMContent:  in.UTMContent,
		ExpiresAt:   in.ExpiresAt,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	err = u.urlMappingRepo.Store(ctx, urlMapping)
	if err != nil {
		if errors.Is(err, entities.ErrURLAlreadyExists) {
			return nil, entities.ErrURLAlreadyExists
		}
		log.Println(err)
		return nil, entities.ErrInternalServer
	}

	if err = u.urlMappingProducer.ProduceURLMappingCreatedEvent(ctx, events.URLCreatedEvent{
		ShortCode:   urlMapping.ShortCode,
		OriginalURL: urlMapping.OriginalURL,
		UTMSource:   urlMapping.UTMSource,
		UTMMedium:   urlMapping.UTMMedium,
		UTMCampaign: urlMapping.UTMCampaign,
		UTMTerm:     urlMapping.UTMTerm,
		UTMContent:  urlMapping.UTMContent,
		ExpiresAt:   urlMapping.ExpiresAt,
	}); err != nil {
		log.Println("Error producing URL mapping created event:", err)
	}

	return &urlMapping, nil
}

func generateShortCode() (string, error) {
	length := 8
	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, length)
	charsetLength := big.NewInt(int64(len(charset)))

	for i := 0; i < length; i++ {
		// Generate a cryptographically secure random number
		randomInt, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", err
		}
		// Use the random number as an index in the charset
		result[i] = charset[randomInt.Int64()]
	}

	return string(result), nil
}

func (u *URLMappingUseCase) GetURLMappingByUUID(ctx context.Context, uuid, ownerUUID string) (*entities.URLMapping, error) {
	url, err := u.urlMappingRepo.GetByUUID(ctx, uuid, ownerUUID)
	if err != nil {
		if errors.Is(err, entities.ErrURLNotFound) {
			return nil, entities.ErrURLNotFound
		}
		log.Println(err)
		return nil, entities.ErrInternalServer
	}
	return url, nil
}

func (u *URLMappingUseCase) GetAllURLMapping(ctx context.Context, ownerUUID string, in pagination.Pagination) (*pagination.PaginatedResponse[entities.URLMapping], error) {
	var urlMappings []entities.URLMapping
	var total int64
	var allowedSortBy = map[string]bool{
		"created_at": true,
		"title":      true,
	}
	if in.Page < 1 {
		in.Page = 1
	}

	if in.PageSize < 1 {
		in.PageSize = 10
	}

	if in.PageSize > 100 {
		in.PageSize = 100
	}

	if in.Sort != "asc" && in.Sort != "desc" {
		in.Sort = "asc"
	}

	sortBy := "created_at" // default
	if allowedSortBy[in.SortBy] {
		sortBy = in.SortBy
	}

	in.SortBy = sortBy

	urlMappings, total, err := u.urlMappingRepo.GetAll(ctx, ownerUUID, in)
	if err != nil {
		log.Println(err)
		return nil, entities.ErrInternalServer
	}

	paginatedResponse := pagination.PaginatedResponse[entities.URLMapping]{
		Items:    urlMappings,
		Total:    total,
		Page:     in.Page,
		PageSize: in.PageSize,
		Sort:     in.Sort,
		SortBy:   in.SortBy,
	}

	return &paginatedResponse, nil
}

func (u *URLMappingUseCase) DeleteURLMapping(ctx context.Context, uuid, ownerUUID string) error {
	err := u.urlMappingRepo.Delete(ctx, uuid, ownerUUID)
	if err != nil {
		if errors.Is(err, entities.ErrURLNotFound) {
			return entities.ErrURLNotFound
		}
		log.Println(err)
		return entities.ErrInternalServer
	}

	if err = u.urlMappingProducer.ProduceURLMappingDeletedEvent(ctx, events.URLDeletedEvent{
		UUID: uuid,
	}); err != nil {
		log.Println("Error producing URL mapping deleted event:", err)
	}

	return nil
}

func (u *URLMappingUseCase) UpdateURLMapping(ctx context.Context, uuid, ownerUUID string, in entities.URLMapping) error {
	err := u.urlMappingRepo.Update(ctx, uuid, ownerUUID, in)
	if err != nil {
		if errors.Is(err, entities.ErrURLNotFound) {
			return entities.ErrURLNotFound
		}
		if errors.Is(err, entities.ErrURLAlreadyExists) {
			return entities.ErrURLAlreadyExists
		}
		log.Println(err)
		return entities.ErrInternalServer
	}
	if err = u.urlMappingProducer.ProduceURLMapingUpdatedEvent(ctx, events.URLUpdatedEvent{
		UUID:        uuid,
		OriginalURL: in.OriginalURL,
		UTMSource:   in.UTMSource,
		UTMMedium:   in.UTMMedium,
		UTMCampaign: in.UTMCampaign,
		UTMTerm:     in.UTMTerm,
		UTMContent:  in.UTMContent,
	}); err != nil {
		log.Println("Error producing URL mapping updated event:", err)
	}
	return nil
}
