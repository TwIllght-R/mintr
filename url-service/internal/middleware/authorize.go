package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authorize(perms []string) gin.HandlerFunc {
	return func(c *gin.Context) {

		userPerms, exists := c.Get("perms")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "No permissions found"})
			return
		}

		userPermList, ok := userPerms.([]string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid permissions format"})
			return
		}

		// Check if user has any of the required permissions
		for _, required := range perms {
			for _, userPerm := range userPermList {
				if strings.EqualFold(required, userPerm) {
					c.Next() // User has permission
					return
				}
			}
		}

		// If no match found
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
	}
}
