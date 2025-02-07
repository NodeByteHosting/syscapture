package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	bearerPrefix        = "Bearer"
	errorParsingHeader  = "Unable to parse 'Authorization' header"
	errorTokenRequired  = "Authorization token required"
	errorInvalidToken   = "Invalid token provided"
)

// AuthRequired is a middleware function that checks for a valid Bearer token in the Authorization header.
func AuthRequired(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(authorizationHeader)
		splittedHeader := strings.Split(authHeader, " ")

		if len(splittedHeader) != 2 || splittedHeader[0] != bearerPrefix {
			c.JSON(401, gin.H{
				"error": errorParsingHeader,
			})
			c.Abort()
			return
		}

		token := splittedHeader[1]
		if token == "" {
			c.JSON(401, gin.H{"error": errorTokenRequired})
			c.Abort()
			return
		} else if token != secret {
			c.JSON(403, gin.H{"error": errorInvalidToken})
			c.Abort()
			return
		}

		c.Next()
	}
}
