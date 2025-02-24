package middleware

import (
	"crypto/subtle"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nodebytehosting/syscapture/internal/config"
)

var (
	ErrMissingHeader = errors.New("missing authorization header")
	ErrInvalidFormat = errors.New("invalid authorization header format")
	ErrTokenRequired = errors.New("authorization token required")
	ErrInvalidToken  = errors.New("invalid token provided")
	ErrTokenExpired  = errors.New("token has expired")
)

// AuthConfig holds configuration for the authentication middleware
type AuthConfig struct {
	Enabled         bool
	Secret          string
	TokenExpiration time.Duration
	SkipPaths       []string
	AllowedHeaders  []string
	RateLimit       int
}

// DefaultAuthConfig returns default authentication configuration
func DefaultAuthConfig() *AuthConfig {
	return &AuthConfig{
		Enabled:         true,
		TokenExpiration: 24 * time.Hour,
		RateLimit:       60,
		AllowedHeaders:  []string{"Authorization", "Content-Type"},
	}
}

// NewAuthConfig creates AuthConfig from application config
func NewAuthConfig(cfg *config.Config) *AuthConfig {
	return &AuthConfig{
		Enabled:         cfg.Security.Auth.Enabled,
		Secret:          cfg.Security.Auth.Secret,
		TokenExpiration: cfg.Security.Auth.TokenExpiry,
		SkipPaths:       cfg.Security.Auth.SkipPaths,
		AllowedHeaders:  cfg.Security.Auth.AllowedHeaders,
		RateLimit:       cfg.Security.Auth.RateLimit.Limit,
	}
}

// AuthRequired is a middleware function that checks for a valid Bearer token
func AuthRequired(config *AuthConfig) gin.HandlerFunc {
	if config == nil {
		config = DefaultAuthConfig()
	}

	return func(c *gin.Context) {
		// Log authentication status for debugging
		c.Set("auth_config", config)

		// If auth is not enabled, skip authentication
		if !config.Enabled {
			c.Set("auth_disabled", true)
			c.Next()
			return
		}

		// Skip authentication for specified paths
		for _, path := range config.SkipPaths {
			if strings.HasPrefix(c.Request.URL.Path, path) {
				c.Set("auth_skipped", true)
				c.Next()
				return
			}
		}

		// Extract and validate token
		token, err := extractToken(c)
		if err != nil {
			handleAuthError(c, err)
			return
		}

		// Check if secret is configured
		if config.Secret == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Authentication secret not configured",
				"code":  http.StatusInternalServerError,
			})
			c.Abort()
			return
		}

		// Constant-time comparison to prevent timing attacks
		if subtle.ConstantTimeCompare([]byte(token), []byte(config.Secret)) != 1 {
			handleAuthError(c, ErrInvalidToken)
			return
		}

		// Store authentication info in context
		c.Set("authenticated", true)
		c.Set("auth_time", time.Now().UTC())

		c.Next()
	}
}

// extractToken extracts the Bearer token from the Authorization header
func extractToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", ErrMissingHeader
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", ErrInvalidFormat
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", ErrTokenRequired
	}

	return token, nil
}

// handleAuthError handles authentication errors with appropriate responses
func handleAuthError(c *gin.Context, err error) {
	var status int
	var message string

	switch err {
	case ErrMissingHeader, ErrInvalidFormat, ErrTokenRequired:
		status = http.StatusUnauthorized
		message = "Authentication required"
	case ErrInvalidToken:
		status = http.StatusForbidden
		message = "Invalid authentication token"
	case ErrTokenExpired:
		status = http.StatusUnauthorized
		message = "Token has expired"
	default:
		status = http.StatusInternalServerError
		message = "Authentication error"
	}

	c.JSON(status, gin.H{
		"error": message,
		"code":  status,
	})
	c.Abort()
}
