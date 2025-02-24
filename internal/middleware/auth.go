package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nodebytehosting/syscapture/internal/config"
)

var (
	ErrMissingAuth      = errors.New("missing authentication")
	ErrInvalidAuth      = errors.New("invalid authentication")
	ErrInsufficientRole = errors.New("insufficient permissions")
)

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleUser   Role = "user"
	RoleViewer Role = "viewer"
)

type APIKey struct {
	Key       string
	Role      Role
	Name      string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type AuthManager struct {
	keys       map[string]APIKey
	config     *AuthConfig
	hashSecret string
}

type AuthConfig struct {
	Enabled        bool
	HashSecret     string
	DefaultRole    Role
	SkipPaths      []string
	AllowedHeaders []string
	RateLimit      RateLimitConfig
}

type RateLimitConfig struct {
	Enabled bool
	Limit   int
	Window  time.Duration
}

func NewAuthManager(cfg *config.SecurityConfig) *AuthManager {
	return &AuthManager{
		keys:       make(map[string]APIKey),
		hashSecret: cfg.Auth.Secret,
		config: &AuthConfig{
			Enabled:        cfg.Auth.Enabled,
			DefaultRole:    Role(cfg.Auth.DefaultRole),
			SkipPaths:      cfg.Auth.SkipPaths,
			AllowedHeaders: cfg.Auth.AllowedHeaders,
			RateLimit: RateLimitConfig{
				Enabled: cfg.Auth.RateLimit.Enabled,
				Limit:   cfg.Auth.RateLimit.Limit,
				Window:  cfg.Auth.RateLimit.Window,
			},
		},
	}
}

func (am *AuthManager) GenerateAPIKey(name string, role Role, expires time.Duration) (string, error) {
	if !am.config.Enabled {
		return "", errors.New("authentication is disabled")
	}

	// Generate unique key using hash
	hash := sha256.New()
	hash.Write([]byte(name))
	hash.Write([]byte(time.Now().String()))
	hash.Write([]byte(am.hashSecret))
	apiKey := hex.EncodeToString(hash.Sum(nil))

	// Store key
	am.keys[apiKey] = APIKey{
		Key:       apiKey,
		Role:      role,
		Name:      name,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(expires),
	}

	return apiKey, nil
}

func (am *AuthManager) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !am.config.Enabled {
			c.Next()
			return
		}

		// Skip authentication for specified paths
		for _, path := range am.config.SkipPaths {
			if strings.HasPrefix(c.Request.URL.Path, path) {
				c.Next()
				return
			}
		}

		// Get API key from header or query
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			apiKey = c.Query("api_key")
		}

		if apiKey == "" {
			handleAuthError(c, ErrMissingAuth)
			return
		}

		// Validate API key
		key, exists := am.keys[apiKey]
		if !exists {
			handleAuthError(c, ErrInvalidAuth)
			return
		}

		// Check expiration
		if time.Now().After(key.ExpiresAt) {
			handleAuthError(c, errors.New("api key expired"))
			return
		}

		// Store auth info in context
		c.Set("api_key", key)
		c.Set("role", key.Role)
		c.Set("user", key.Name)

		c.Next()
	}
}

// RequireRole middleware checks if authenticated user has required role
func (am *AuthManager) RequireRole(role Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			handleAuthError(c, ErrInsufficientRole)
			return
		}

		if userRole != role && userRole != RoleAdmin {
			handleAuthError(c, ErrInsufficientRole)
			return
		}

		c.Next()
	}
}

func handleAuthError(c *gin.Context, err error) {
	status := http.StatusUnauthorized
	message := "Authentication failed"

	switch err {
	case ErrMissingAuth:
		message = "Authentication required"
	case ErrInvalidAuth:
		message = "Invalid API key"
	case ErrInsufficientRole:
		status = http.StatusForbidden
		message = "Insufficient permissions"
	}

	c.JSON(status, gin.H{
		"error": message,
		"code":  status,
	})
	c.Abort()
}
