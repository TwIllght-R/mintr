package interfaces

import (
	"context"
	"redirect-service/domain/entities"
)

type RedirectURLUsecase interface {
	Redirect(ctx context.Context, in entities.URLVisit) (*string, error)
	CreateRedirectURL(ctx context.Context, in entities.RedirectURL) (*entities.RedirectURL, error)
	UpdateRedirectURL(ctx context.Context, uuid string, in entities.RedirectURL) (*entities.RedirectURL, error)
	DeleteRedirectURL(ctx context.Context, uuid string) error
}

type RedirectURLCache interface {
	StoreMappingURL(ctx context.Context, shortCode, originalURLWithUTM string) error
	GetMappingURL(ctx context.Context, shortCode string) (*string, error)
}

type RedirectURLRepository interface {
	Store(ctx context.Context, in entities.RedirectURL) (*entities.RedirectURL, error)
	GetByShortCode(ctx context.Context, shortCode string) (*entities.RedirectURL, error)
	GetByUUID(ctx context.Context, uuid string) (*entities.RedirectURL, error)
	Update(ctx context.Context, uuid string, in entities.RedirectURL) (*entities.RedirectURL, error)
	Delete(ctx context.Context, uuid string) error
}
