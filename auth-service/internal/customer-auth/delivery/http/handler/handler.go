package handler

import (
	"auth-service/domain/entities"
	"auth-service/domain/interfaces"
	"auth-service/internal/customer-auth/delivery/http/dto"
	"auth-service/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type customerAuthHandler struct {
	customerAuthUseCase interfaces.CustomerAuthUseCase
	jwtSecret           []byte
}

func NewCustomerAuthHandler(g *gin.Engine, customerAuthUseCase interfaces.CustomerAuthUseCase, jwtSecret []byte) {
	handler := &customerAuthHandler{
		customerAuthUseCase: customerAuthUseCase,
		jwtSecret:           jwtSecret,
	}

	v1 := g.Group("/v1")

	authGroups := v1.Group("/auth")
	authGroups.POST("/sign-ups", handler.SignUp)
	authGroups.GET("/email-verifications", handler.VerifyEmailVerifications)
	authGroups.POST("/logins", handler.Login)
	authGroups.POST("/logouts", handler.Logout)

}

func mapStatusCode(err error) int {
	switch err {
	case entities.ErrUserNotFound:
		return http.StatusNotFound
	case entities.ErrEmailAlreadyExists:
		return http.StatusConflict
	case entities.ErrEmailNotVerified:
		return http.StatusForbidden
	case entities.ErrEmailAlreadyVerified:
		return http.StatusConflict
	case entities.ErrInternalServer:
		return http.StatusInternalServerError
	case entities.ErrInvalidCredentials:
		return http.StatusUnauthorized
	case entities.ErrPermissionNotFound:
		return http.StatusNotFound
	case entities.ErrRoleNotFound:
		return http.StatusNotFound
	case entities.ErrRoleAlreadyExists:
		return http.StatusConflict
	case entities.ErrPermissionAlreadyExists:
		return http.StatusConflict
	case entities.ErrRolePermissionNotFound:
		return http.StatusNotFound
	case entities.ErrInvalidToken:
		return http.StatusUnauthorized
	case entities.ErrGenerateToken:
		return http.StatusInternalServerError
	case entities.ErrRevokedToken:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

func (h *customerAuthHandler) SignUp(c *gin.Context) {
	var req dto.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	err := h.customerAuthUseCase.SignUp(c.Request.Context(), req.ToEntity())
	if err != nil {
		c.JSON(mapStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "success",
	})
}

func (h *customerAuthHandler) VerifyEmailVerifications(c *gin.Context) {
	token, ok := c.GetQuery("token")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token is required"})
		return
	}

	err := h.customerAuthUseCase.VerifyEmail(c.Request.Context(), token)
	if err != nil {
		c.JSON(mapStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "success",
	})

}

func (h *customerAuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	resp, err := h.customerAuthUseCase.Login(c.Request.Context(), req.ToEntity())
	if err != nil {
		c.JSON(mapStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"data":    dto.ToLoginResponse(*resp),
		"message": "success",
	})
}

func (h *customerAuthHandler) Logout(c *gin.Context) {
	var req dto.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Get the JTI from the access token in the request header
	jwtToken := c.Request.Header.Get("Authorization")
	if jwtToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header is required"})
		return
	}

	jwtToken = jwtToken[len("Bearer "):]

	if jwtToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "access token is required"})
		return
	}

	claimsToken, err := utils.ExtractClaims(jwtToken, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		return
	}

	jti, ok := claimsToken["jti"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		return
	}

	exp, ok := claimsToken["exp"].(float64)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		return
	}

	err = h.customerAuthUseCase.Logout(c.Request.Context(), jti, req.RefreshToken, time.Unix(int64(exp), 0))
	if err != nil {
		c.JSON(mapStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "success",
	})
}
