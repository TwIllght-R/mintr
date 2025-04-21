package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewURLMappingHandler(g *gin.Engine) {

	g.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

}
