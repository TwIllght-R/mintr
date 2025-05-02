package handler

import (
	"analytics-service/domain/entities"
	"analytics-service/domain/interfaces"
	"analytics-service/internal/analytics/delivery/http/dto"
	"analytics-service/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type analyticsHandler struct {
	analyticsUseCase interfaces.AnalyticsUseCase
	jwtSecret        []byte
}

func NewAnalyticsHandler(g *gin.Engine, jwtSecret []byte, analyticsUseCase interfaces.AnalyticsUseCase) {
	handler := &analyticsHandler{
		jwtSecret:        jwtSecret,
		analyticsUseCase: analyticsUseCase,
	}

	v1 := g.Group("/v1")

	url := v1.Group("/analytics", middleware.JWTAuth(jwtSecret))

	url.Use(middleware.Authorize([]string{"analytics:read"}))
	url.GET("/:url-uuid", handler.GetClickedDetails)

}

func mapStatusCode(err error) int {
	switch err {
	case entities.ErrInternalServer:
		return http.StatusInternalServerError
	default:
		return http.StatusBadRequest
	}
}

func (h *analyticsHandler) GetClickedDetails(c *gin.Context) {
	urlUUID := c.Param("url-uuid")
	if urlUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url-uui is required"})
		return
	}
	resp, err := h.analyticsUseCase.GetClickedDetails(c.Request.Context(), urlUUID, c.GetString("sub"))
	if err != nil {
		c.JSON(mapStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    dto.ToURLClickedResponse(resp),
		"message": "success",
	})
}
