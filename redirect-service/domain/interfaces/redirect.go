package interfaces

import (
	"context"
	"redirect-service/domain/entities"
)

type RedirectURLUsecase interface {
	Redirect(ctx context.Context, in entities.RedirectURL) (*string, error)
}

type RedirectURLCache interface {
	Cache(ctx context.Context, shortCode, originalURL string) error
	Get(ctx context.Context, key string) (*string, error)
}

type RedirectURLRepository interface {
	GetByShortCode(ctx context.Context, shortCode string) (*entities.RedirectURL, error)
}
