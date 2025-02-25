package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/Cassandra-Labs-Foundation/core/internal/service/auth"
)

// Handler provides authentication HTTP handlers
type Handler struct {
	service auth.Service
}

// NewHandler creates a new authentication handler
func NewHandler(service auth.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// TokenResponse represents the token response
type TokenResponse struct {
	Token string `json:"token"`
}

// Login handles the login request
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	token, err := h.service.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	
	c.JSON(http.StatusOK, TokenResponse{Token: token})
}

// RefreshToken handles the token refresh request
func (h *Handler) RefreshToken(c *gin.Context) {
	// Get the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}
	
	// Check that the Authorization header has the format "Bearer {token}"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
		return
	}
	
	// Extract the token
	tokenString := parts[1]
	
	// Refresh the token
	newToken, err := h.service.RefreshToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}
	
	c.JSON(http.StatusOK, TokenResponse{Token: newToken})
}

// ValidateToken handles the token validation request
func (h *Handler) ValidateToken(c *gin.Context) {
	// This endpoint is just for testing purposes
	// In a real application, you would probably not expose this endpoint
	
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	
	c.JSON(http.StatusOK, gin.H{
		"userID": userID,
		"role":   role,
	})
}