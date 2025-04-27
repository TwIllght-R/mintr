package dto

import "auth-service/domain/entities"

type SignUpRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (d *SignUpRequest) ToEntity() entities.CustomerAuth {
	return entities.CustomerAuth{
		Email:    d.Email,
		Password: d.Password,
	}
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (d *LoginRequest) ToEntity() entities.CustomerAuth {
	return entities.CustomerAuth{
		Email:    d.Email,
		Password: d.Password,
	}
}

type LoginResponse struct {
	AccessToken           string `json:"access_token"`
	RefreshToken          string `json:"refresh_token"`
	AccessTokenExpiresAt  string `json:"access_token_expires_at"`
	RefreshTokenExpiresAt string `json:"refresh_token_expires_at"`
}

func ToLoginResponse(in entities.Token) LoginResponse {
	return LoginResponse{
		AccessToken:           in.AccessToken,
		RefreshToken:          in.RefreshToken,
		AccessTokenExpiresAt:  in.AccessTokenExpiresAt.Format("2006-01-02 15:04:05"),
		RefreshTokenExpiresAt: in.RefreshTokenExpiresAt.Format("2006-01-02 15:04:05"),
	}
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
