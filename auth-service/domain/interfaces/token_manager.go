package interfaces

import "time"

type TokenManager interface {
	// ─── Access / Refresh JWTs ───
	GenerateAccessToken(userUUID, role string, perms []string) (token string, expiresAt time.Time, err error)
	//ValidateAccessToken(token string) (userID, role, jti string, perms []string, err error)

	// ─── Refresh Token ───
	GenerateRefreshToken(userUUID string) (token string, expiresAt time.Time, err error)
	ValidateRefreshToken(userUUID, token string) (err error)

	// // ─── Revocation / Blacklist ───
	RevokeAccessToken(jti string, expiresAt time.Time) error
	RevokeRefreshToken(token string) error

	// // ─── Helpers ───
	// ExtractClaims(token string) (claims map[string]interface{}, err error)
}

type RevocationStore interface {
	IsRevoked(jti string) bool
	RevokeJTI(jti string, ttl time.Duration) error
	StoreRefreshToken(token, userUUID string, ttl time.Duration) error
	DeleteRefreshToken(token string) error
	GetRefreshToken(token string) (string, error)
}

// type EmailTokenManager interface {
// 	// ─── Email Verification ───
// 	GenerateEmailVerificationToken(userUUID, email string) (token string, expiresAt time.Time, err error)
// 	ValidateEmailVerificationToken(token string) (userID, email string, err error)
// 	// ─── Password Reset ───
// 	// GeneratePasswordResetToken(userID, email string) (token string, expiresAt time.Time, err error)
// 	// ValidatePasswordResetToken(token string) (userID, email string, err error)
// }
