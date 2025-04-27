package usecase

import (
	"auth-service/domain/entities"
	"auth-service/domain/events"
	"auth-service/domain/interfaces"
	"auth-service/pkg/utils"
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
)

type customerAuthUseCase struct {
	customerAuthRepo          interfaces.CustomerAuthRepository
	customerAuthCache         interfaces.CustomerAuthCache
	customerAuthEventProducer interfaces.CustomerAuthEventProducer
	tokenManager              interfaces.TokenManager
	permissionUseCase         interfaces.PermissionUseCase
}

func NewCustomerAuthUseCase(
	customerAuthRepo interfaces.CustomerAuthRepository,
	customerAuthCache interfaces.CustomerAuthCache,
	customerAuthEventProducer interfaces.CustomerAuthEventProducer,
	tokenManager interfaces.TokenManager,
	permissionUseCase interfaces.PermissionUseCase,
) interfaces.CustomerAuthUseCase {
	return &customerAuthUseCase{
		customerAuthRepo:          customerAuthRepo,
		customerAuthCache:         customerAuthCache,
		customerAuthEventProducer: customerAuthEventProducer,
		tokenManager:              tokenManager,
		permissionUseCase:         permissionUseCase,
	}
}

func (u *customerAuthUseCase) SignUp(ctx context.Context, in entities.CustomerAuth) error {
	const RoleUUID = "00000000-0000-0000-0000-000000000001"
	//check if email already exists
	existing, err := u.customerAuthRepo.GetByEmail(ctx, in.Email)
	if err != nil && err != entities.ErrUserNotFound {
		log.Println(err)
		return entities.ErrInternalServer
	}
	if existing != nil && existing.IsVerified {
		return entities.ErrEmailAlreadyExists
	}

	hashPassword, err := utils.HashPassword(in.Password)
	if err != nil {
		log.Println("Error hashing password:", err)
		return entities.ErrInternalServer
	}

	newCustomer := entities.CustomerAuth{
		UUID:       uuid.New().String(),
		Email:      in.Email,
		Password:   hashPassword,
		IsVerified: false,
		RoleUUID:   RoleUUID,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
	err = u.customerAuthRepo.Store(ctx, newCustomer)
	if err != nil {
		if errors.Is(err, entities.ErrEmailAlreadyExists) {
			return entities.ErrEmailAlreadyExists
		}
		log.Println(err)
		return entities.ErrInternalServer
	}

	token := uuid.New().String()
	err = u.customerAuthCache.StoreEmailVerificationToken(ctx, newCustomer.Email, token)
	if err != nil {
		log.Println("Error caching token:", err)
		return entities.ErrInternalServer
	}

	//produce email verification event
	if err := u.customerAuthEventProducer.ProduceEmailVerificationEvent(ctx, events.EmailVerificationEvent{
		UserUUID: newCustomer.UUID,
		Email:    newCustomer.Email,
		Token:    token,
	}); err != nil {
		log.Println("Error producing email verification event:", err)
	}

	return nil
}

func (u *customerAuthUseCase) Login(ctx context.Context, in entities.CustomerAuth) (*entities.Token, error) {
	// Check if the user exists in the database
	existing, err := u.customerAuthRepo.GetByEmail(ctx, in.Email)
	if err != nil {
		if errors.Is(err, entities.ErrUserNotFound) {
			return nil, entities.ErrInvalidCredentials
		}
		log.Println(err)
		return nil, entities.ErrInternalServer
	}

	// Verify the password
	if !utils.CheckPasswordHash(in.Password, existing.Password) {
		return nil, entities.ErrInvalidCredentials
	}

	// Check if the user is verified
	if !existing.IsVerified {
		return nil, entities.ErrEmailNotVerified
	}

	//get permissions
	role, perms, err := u.permissionUseCase.GetRoleAndPermissionByRoleUUID(ctx, existing.RoleUUID)
	if err != nil {
		log.Println(err)
		return nil, entities.ErrInternalServer
	}

	// Generate access and refresh tokens
	accessToken, accTkExp, err := u.tokenManager.GenerateAccessToken(existing.UUID, role.Name, func() []string {
		var permissions []string
		for _, perm := range perms {
			permissions = append(permissions, perm.Name)
		}
		return permissions
	}())

	if err != nil {
		log.Println("Error generating access token:", err)
		return nil, entities.ErrGenerateToken
	}
	refreshToken, refTkExp, err := u.tokenManager.GenerateRefreshToken(existing.UUID)
	if err != nil {
		log.Println("Error generating refresh token:", err)
		return nil, entities.ErrGenerateToken
	}

	return &entities.Token{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accTkExp,
		RefreshTokenExpiresAt: refTkExp,
	}, nil
}

func (u *customerAuthUseCase) Logout(ctx context.Context, jti, refreshToken string, accessTokenExpires time.Time) error {
	var err error
	err = u.tokenManager.RevokeAccessToken(jti, accessTokenExpires)
	if err != nil {
		log.Println("Error revoking access token:", err)
		return entities.ErrRevokedToken
	}
	err = u.tokenManager.RevokeRefreshToken(refreshToken)
	if err != nil {
		log.Println("Error revoking refresh token:", err)
		return entities.ErrRevokedToken
	}
	return nil

}

func (u *customerAuthUseCase) VerifyEmail(ctx context.Context, token string) error {
	email, err := u.customerAuthCache.VerifyEmailVerificationToken(ctx, token)
	if err != nil {
		log.Println("Error verifying email token:", err)
		return entities.ErrInvalidToken
	}

	auth, err := u.customerAuthRepo.GetByEmail(ctx, *email)
	if err != nil {
		if errors.Is(err, entities.ErrUserNotFound) {
			return entities.ErrUserNotFound
		}
		log.Println(err)
		return entities.ErrInternalServer
	}
	auth.IsVerified = true
	auth.UpdatedAt = time.Now().UTC()

	err = u.customerAuthRepo.Update(ctx, auth.UUID, *auth)
	if err != nil {
		if errors.Is(err, entities.ErrUserNotFound) {
			return entities.ErrUserNotFound
		}
		log.Println(err)
		return entities.ErrInternalServer
	}

	return nil
}

func (u *customerAuthUseCase) ResendVerificationEmail(ctx context.Context, email string) error {
	existing, err := u.customerAuthRepo.GetByEmail(ctx, email)
	if err != nil {
		return entities.ErrInvalidCredentials
	}

	if existing.IsVerified {
		return entities.ErrEmailAlreadyVerified
	}

	token := uuid.New().String()
	err = u.customerAuthCache.StoreEmailVerificationToken(ctx, existing.Email, token)
	if err != nil {
		log.Println("Error caching token:", err)
		return entities.ErrInternalServer
	}

	// Produce email verification event
	if err := u.customerAuthEventProducer.ProduceEmailVerificationEvent(ctx, events.EmailVerificationEvent{
		UserUUID: existing.UUID,
		Email:    existing.Email,
		Token:    token,
	}); err != nil {
		log.Println("Error producing email verification event:", err)

	}

	return nil
}

func (u *customerAuthUseCase) RefreshToken(ctx context.Context, userUUID, refreshToken string) (*entities.Token, error) {
	// Verify the refresh token
	if err := u.tokenManager.ValidateRefreshToken(userUUID, refreshToken); err != nil {
		log.Println("Error verifying refresh token:", err)
		return nil, entities.ErrInvalidToken
	}

	customer, err := u.customerAuthRepo.GetByUUID(ctx, userUUID)
	if err != nil {
		if errors.Is(err, entities.ErrUserNotFound) {
			return nil, entities.ErrUserNotFound
		}
		log.Println(err)
		return nil, entities.ErrInternalServer
	}

	//get permissions
	role, perms, err := u.permissionUseCase.GetRoleAndPermissionByRoleUUID(ctx, customer.RoleUUID)
	if err != nil {
		log.Println(err)
		return nil, entities.ErrInternalServer
	}

	// Generate new access and refresh tokens
	accessToken, accTkExp, err := u.tokenManager.GenerateAccessToken(userUUID, role.Name, func() []string {
		var permissions []string
		for _, perm := range perms {
			permissions = append(permissions, perm.Name)
		}
		return permissions
	}())
	if err != nil {
		log.Println("Error generating access token:", err)
		return nil, entities.ErrGenerateToken
	}
	newRefreshToken, refTkExp, err := u.tokenManager.GenerateRefreshToken(userUUID)
	if err != nil {
		log.Println("Error generating refresh token:", err)
		return nil, entities.ErrGenerateToken
	}

	return &entities.Token{
		AccessToken:           accessToken,
		RefreshToken:          newRefreshToken,
		AccessTokenExpiresAt:  accTkExp,
		RefreshTokenExpiresAt: refTkExp,
	}, nil
}
