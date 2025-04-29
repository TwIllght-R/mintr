package handler

import (
	"net/http"
	"url-service/domain/entities"
	"url-service/domain/interfaces"

	"url-service/internal/middleware"
	"url-service/internal/url-mapping/delivery/http/dto"

	"github.com/gin-gonic/gin"
)

type urlMappingHandler struct {
	urlMappingUseCase interfaces.URLMappingUseCase
	jwtSecret         []byte
}

func NewURLMappingHandler(g *gin.Engine, jwtSecret []byte, urlMappingUseCase interfaces.URLMappingUseCase) {
	handler := &urlMappingHandler{
		jwtSecret:         jwtSecret,
		urlMappingUseCase: urlMappingUseCase,
	}

	v1 := g.Group("/v1")

	url := v1.Group("/url", middleware.JWTAuth(jwtSecret))

	url.Use(middleware.Authorize([]string{"url:read", "url:write"}))
	url.POST("", handler.GenerateShortCode)
	url.GET("", handler.GetAllURLMapping)
	url.GET("/:uuid", handler.GetURLMappingByUUID)
	url.PATCH("/:uuid", handler.UpdateURLMapping)
	url.DELETE("/:uuid", handler.DeleteURLMapping)

}

func mapStatusCode(err error) int {
	switch err {
	case entities.ErrURLAlreadyExists:
		return http.StatusConflict
	case entities.ErrURLNotFound:
		return http.StatusNotFound
	case entities.ErrInternalServer:
		return http.StatusInternalServerError
	default:
		return http.StatusBadRequest
	}
}

func (h *urlMappingHandler) GenerateShortCode(c *gin.Context) {
	var req dto.GenerateShortCodeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	resp, err := h.urlMappingUseCase.GenerateShortCode(c.Request.Context(), c.GetString("sub"), req.ToEntity())
	if err != nil {
		c.JSON(mapStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"data":    dto.ToURLMappingResponse(*resp),
		"message": "success",
	})
}

func (h *urlMappingHandler) GetAllURLMapping(c *gin.Context) {
	var req dto.GetAllMappingURLRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	resp, err := h.urlMappingUseCase.GetAllURLMapping(c.Request.Context(), c.GetString("sub"), req.ToEntity())
	if err != nil {
		c.JSON(mapStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    dto.ToResponsePaginated(resp.Items, resp.Page, resp.PageSize, resp.Sort, resp.SortBy, resp.Total),
		"message": "success",
	})
}

func (h *urlMappingHandler) GetURLMappingByUUID(c *gin.Context) {
	uuid := c.Param("uuid")
	if uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uuid is required"})
		return
	}
	resp, err := h.urlMappingUseCase.GetURLMappingByUUID(c.Request.Context(), uuid, c.GetString("sub"))
	if err != nil {
		c.JSON(mapStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    dto.ToURLMappingResponse(*resp),
		"message": "success",
	})

}

func (h *urlMappingHandler) UpdateURLMapping(c *gin.Context) {
	uuid := c.Param("uuid")
	if uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uuid is required"})
		return
	}
	var req dto.UpdateURLMappingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	err := h.urlMappingUseCase.UpdateURLMapping(c.Request.Context(), uuid, c.GetString("sub"), req.ToEntity())
	if err != nil {
		c.JSON(mapStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "success",
	})
}

func (h *urlMappingHandler) DeleteURLMapping(c *gin.Context) {
	uuid := c.Param("uuid")
	if uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uuid is required"})
		return
	}
	err := h.urlMappingUseCase.DeleteURLMapping(c.Request.Context(), uuid, c.GetString("sub"))
	if err != nil {
		c.JSON(mapStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
