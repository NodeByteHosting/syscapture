package middleware

import (
	"strings"
	"github.com/gin-gonic/gin"
)

// AuthRequired is a middleware function that checks for a valid Bearer token in the Authorization header.
func AuthRequired(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		splittedHeader := strings.Split(authHeader, " ")

		// Check if the Authorization header is properly formatted
		if len(splittedHeader) != 2 || splittedHeader[0] != "Bearer" {
			c.JSON(401, gin.H{
				"error": "Unable to parse 'Authorization' header",
			})
			c.Abort()
			return
		}

		token := splittedHeader[1]

		// Check if the token is provided
		if token == "" {
			c.JSON(401, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		// Check if the token matches the secret
		if token != secret {
			c.JSON(403, gin.H{"error": "Invalid token provided"})
			c.Abort()
			return
		}

		c.Next()
	}
}
