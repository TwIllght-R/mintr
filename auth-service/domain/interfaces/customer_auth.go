package interfaces

import (
	"auth-service/domain/entities"
	"auth-service/domain/events"
	"context"
	"time"
)

type CustomerAuthUseCase interface {
	SignUp(ctx context.Context, in entities.CustomerAuth) error
	Login(ctx context.Context, in entities.CustomerAuth) (*entities.Token, error)
	Logout(ctx context.Context, jti, refreshToken string, accessTokenExpires time.Time) error
	VerifyEmail(ctx context.Context, token string) error
	ResendVerificationEmail(ctx context.Context, email string) error
	RefreshToken(ctx context.Context, userUUID, refreshToken string) (*entities.Token, error)
}

type CustomerAuthRepository interface {
	Store(ctx context.Context, in entities.CustomerAuth) error
	Update(ctx context.Context, uuid string, in entities.CustomerAuth) error
	GetByUUID(ctx context.Context, uuid string) (*entities.CustomerAuth, error)
	GetByEmail(ctx context.Context, email string) (*entities.CustomerAuth, error)
}

type CustomerAuthCache interface {
	StoreEmailVerificationToken(ctx context.Context, email, token string) error
	VerifyEmailVerificationToken(ctx context.Context, token string) (*string, error)
}

type CustomerAuthEventProducer interface {
	ProduceEmailVerificationEvent(ctx context.Context, in events.EmailVerificationEvent) error
}
