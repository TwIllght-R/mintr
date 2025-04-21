package usecase

import (
	"url-service/domain/entities"
	"url-service/domain/interfaces"
	"url-service/domain/pagination"
	"context"
	"crypto/rand"
	"log"
	"math/big"
	"time"

	"github.com/google/uuid"
)

type URLMappingUseCase struct {
	urlMappingRepo interfaces.URLMappingRepository
}

func NewURLMappingUseCase(urlMappingRepo interfaces.URLMappingRepository) interfaces.URLMappingUseCase {
	return &URLMappingUseCase{urlMappingRepo: urlMappingRepo}
}

func (u *URLMappingUseCase) GenerateShortCode(ctx context.Context, in entities.URLMapping) (*entities.URLMapping, error) {
	// Generate a unique short code
	shortCode, err := generateShortCode()
	if err != nil {
		log.Println("Error generating short code:", err)
		return nil, entities.ErrInternalServer
	}

	
	urlMapping := entities.URLMapping{
		UUID:        uuid.New().String(),
		ShortCode:   shortCode,
		OriginalURL: in.OriginalURL,
		Title:       in.Title,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	
	err = u.urlMappingRepo.Store(ctx, urlMapping)
	if err != nil {
		return nil, err
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

func (u *URLMappingUseCase) GetURLMappingByUUID(ctx context.Context, uuid string) (*entities.URLMapping, error) {
	return u.urlMappingRepo.GetByUUID(ctx, uuid)
}

func (u *URLMappingUseCase) GetAllURLMapping(ctx context.Context, in pagination.Pagination) (*pagination.PaginatedResponse[entities.URLMapping], error) {
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

	
	urlMappings, total, err := u.urlMappingRepo.GetAll(ctx, in)
	if err != nil {
		return nil, err
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

func (u *URLMappingUseCase) DeleteURLMapping(ctx context.Context, uuid string) error {
	return u.urlMappingRepo.Delete(ctx, uuid)
}

func (u *URLMappingUseCase) UpdateURLMapping(ctx context.Context, uuid string, in entities.URLMapping) error {
	return u.urlMappingRepo.Update(ctx, uuid, in)
}
