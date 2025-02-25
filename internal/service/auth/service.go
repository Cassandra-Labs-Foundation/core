package auth

import (
	"errors"

	"github.com/Cassandra-Labs-Foundation/core/pkg/jwt"
)

// Service provides authentication business logic
type Service interface {
	Login(username, password string) (string, error)
	RefreshToken(tokenString string) (string, error)
	ValidateToken(tokenString string) (string, string, error)
}

type service struct {
	jwtService jwt.Service
	// In a real application, you would have a repository to store users and their credentials
	// userRepo repository.UserRepository
}

// NewService creates a new authentication service
func NewService(jwtService jwt.Service) Service {
	return &service{
		jwtService: jwtService,
	}
}

// Login authenticates a user and returns a JWT token if successful
func (s *service) Login(username, password string) (string, error) {
	// In a real application, you would look up the user in a database
	// and verify the password. For this example, we'll use hardcoded credentials.
	
	// This is just for demo purposes. In a real app, NEVER store passwords in code.
	if username == "admin" && password == "password" {
		// Generate a token for the admin user
		return s.jwtService.GenerateToken("admin-user-id", "admin")
	}
	
	if username == "user" && password == "password" {
		// Generate a token for a regular user
		return s.jwtService.GenerateToken("regular-user-id", "user")
	}
	
	return "", errors.New("invalid credentials")
}

// RefreshToken validates an existing token and returns a new one
func (s *service) RefreshToken(tokenString string) (string, error) {
	// Validate the current token
	claims, err := s.jwtService.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}
	
	// Generate a new token with the same claims
	return s.jwtService.GenerateToken(claims.UserID, claims.Role)
}

// ValidateToken validates a token and returns the user ID and role
func (s *service) ValidateToken(tokenString string) (string, string, error) {
	claims, err := s.jwtService.ValidateToken(tokenString)
	if err != nil {
		return "", "", err
	}
	
	return claims.UserID, claims.Role, nil
}