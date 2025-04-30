package usecase

import (
	"context"
	"errors"
	"log"
	"redirect-service/domain/entities"
	"redirect-service/domain/interfaces"
	"strings"
	"time"
)

type redirectURLUseCase struct {
	redirectURLRepo  interfaces.RedirectURLRepository
	redirectURLCache interfaces.RedirectURLCache
}

func NewRedirectURLUseCase(repo interfaces.RedirectURLRepository, cache interfaces.RedirectURLCache) interfaces.RedirectURLUsecase {
	return &redirectURLUseCase{
		redirectURLRepo:  repo,
		redirectURLCache: cache,
	}
}

func (u *redirectURLUseCase) Redirect(ctx context.Context, in entities.URLVisit) (*string, error) {
	cachedURL, err := u.redirectURLCache.GetMappingURL(ctx, in.ShortCode)
	if err == nil {
		return cachedURL, nil
	}
	if !errors.Is(err, entities.ErrURLNotFound) {
		return nil, entities.ErrInternalServer
	}

	url, err := u.redirectURLRepo.GetByShortCode(ctx, in.ShortCode)
	if err != nil {
		if errors.Is(err, entities.ErrURLNotFound) {
			return nil, entities.ErrURLNotFound
		}
		log.Println(err)
		return nil, entities.ErrInternalServer
	}
	utm := map[string]string{
		"utm_source":   url.UTMSource,
		"utm_medium":   url.UTMMedium,
		"utm_campaign": url.UTMCampaign,
		"utm_term":     url.UTMTerm,
		"utm_content":  url.UTMContent,
	}

	originWithUTM := url.OriginalURL
	first := true
	for key, value := range utm {
		if value != "" {
			if first && !strings.Contains(originWithUTM, "?") {
				originWithUTM += "?" + key + "=" + value
				first = false
			} else {
				originWithUTM += "&" + key + "=" + value
			}
		}
	}
	if err := u.redirectURLCache.StoreMappingURL(ctx, url.ShortCode, originWithUTM); err != nil {
		log.Println(err)
		return nil, entities.ErrInternalServer
	}

	return &originWithUTM, nil
}

func (u *redirectURLUseCase) CreateRedirectURL(ctx context.Context, in entities.RedirectURL) (*entities.RedirectURL, error) {

	in.CreatedAt = time.Now().UTC()
	in.UpdatedAt = time.Now().UTC()

	createdURL, err := u.redirectURLRepo.Store(ctx, in)
	if err != nil {
		log.Println(err)
		return nil, entities.ErrInternalServer
	}
	return createdURL, nil
}

func (u *redirectURLUseCase) UpdateRedirectURL(ctx context.Context, uuid string, in entities.RedirectURL) (*entities.RedirectURL, error) {

	in.UpdatedAt = time.Now().UTC()

	updatedURL, err := u.redirectURLRepo.Update(ctx, uuid, in)
	if err != nil {
		if errors.Is(err, entities.ErrURLNotFound) {
			return nil, entities.ErrURLNotFound
		}
		log.Println(err)
		return nil, entities.ErrInternalServer
	}
	return updatedURL, nil
}

func (u *redirectURLUseCase) DeleteRedirectURL(ctx context.Context, uuid string) error {
	err := u.redirectURLRepo.Delete(ctx, uuid)
	if err != nil {
		if errors.Is(err, entities.ErrURLNotFound) {
			return entities.ErrURLNotFound
		}
		log.Println(err)
		return entities.ErrInternalServer
	}
	return nil
}
