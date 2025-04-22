package usecase

import (
	"context"
	"redirect-service/domain/entities"
	"redirect-service/domain/interfaces"
)

type redirectURLUsecase struct {
	redirectURLRepo  interfaces.RedirectURLRepository
	redirectURLCache interfaces.RedirectURLCache
}

func NewRedirectURLUsecase(repo interfaces.RedirectURLRepository, cache interfaces.RedirectURLCache) interfaces.RedirectURLUsecase {
	return &redirectURLUsecase{
		redirectURLRepo:  repo,
		redirectURLCache: cache,
	}
}

func (u *redirectURLUsecase) Redirect(ctx context.Context, in entities.RedirectURL) (*string, error) {
	cachedURL, err := u.redirectURLCache.Get(ctx, in.ShortCode)
	if err == nil {
		return cachedURL, nil
	}
	if err != entities.ErrURLNotFound {
		return nil, err
	}

	url, err := u.redirectURLRepo.GetByShortCode(ctx, in.ShortCode)
	if err != nil {
		return nil, err
	}

	if err := u.redirectURLCache.Cache(ctx, url.ShortCode, url.OriginalURL); err != nil {
		return nil, err
	}

	return &url.OriginalURL, nil
}
