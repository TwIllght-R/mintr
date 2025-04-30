package handler

import (
	"net/http"
	"redirect-service/domain/entities"
	"redirect-service/domain/interfaces"
	"regexp"

	"github.com/gin-gonic/gin"
)

type redirectURLHandler struct {
	redirectURLUsecase interfaces.RedirectURLUsecase
}

func NewRedirectURLHandler(g *gin.Engine, redirectURLUsecase interfaces.RedirectURLUsecase) {
	handler := &redirectURLHandler{
		redirectURLUsecase: redirectURLUsecase,
	}

	g.GET("/:short-code", handler.RedirectURL)

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

func (h *redirectURLHandler) RedirectURL(c *gin.Context) {
	shortCode := c.Param("short-code")

	if !isValidShortCode(shortCode) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid short code"})
		return
	}

	originalURL, err := h.redirectURLUsecase.Redirect(c.Request.Context(), entities.URLVisit{
		ShortCode: shortCode,
		UserAgent: c.Request.UserAgent(),
		IP:        c.ClientIP(),
	})
	if err != nil {
		c.JSON(mapStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.Redirect(302, *originalURL)
}

func isValidShortCode(shortCode string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9]{1,10}$`)
	return re.MatchString(shortCode)
}
