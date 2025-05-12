package token_manager

import (
	"auth-service/domain/interfaces"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type tokenManager struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
	revocation interfaces.RevocationStore
}

func NewTokenManager(secret []byte, accessTTL, refreshTTL time.Duration, revocation interfaces.RevocationStore) interfaces.TokenManager {
	return &tokenManager{
		secret:     secret,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
		revocation: revocation,
	}
}

func (tm *tokenManager) GenerateAccessToken(userUUID, role string, perms []string) (accessToken string, expiresAt time.Time, err error) {
	exp := time.Now().Add(tm.accessTTL)
	claims := jwt.MapClaims{
		"sub":  userUUID,
		"role": role,
		"perm": perms,
		"exp":  exp.Unix(),
		"jti":  uuid.NewString(),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokStr, err := jwtToken.SignedString(tm.secret)
	return tokStr, exp, err
}

func (tm *tokenManager) GenerateRefreshToken(userUUID string) (refreshToken string, expiresAt time.Time, err error) {
	token := uuid.NewString()
	expiresAt = time.Now().Add(tm.refreshTTL)

	if err := tm.revocation.StoreRefreshToken(token, userUUID, tm.refreshTTL); err != nil {
		return "", time.Time{}, err
	}

	return token, expiresAt, nil
}

func (tm *tokenManager) RevokeAccessToken(jti string, expiresAt time.Time) error {
	return tm.revocation.RevokeJTI(jti, time.Until(expiresAt))
}

func (tm *tokenManager) RevokeRefreshToken(tok string) error {
	return tm.revocation.DeleteRefreshToken(tok)
}

func (tm *tokenManager) ValidateRefreshToken(userUUID, token string) error {
	val, err := tm.revocation.GetRefreshToken(token)
	if err != nil {
		return err
	}
	if val != userUUID {
		return errors.New("invalid refresh token")
	}
	return nil
}
