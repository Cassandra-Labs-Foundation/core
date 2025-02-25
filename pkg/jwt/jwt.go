package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Custom claims structure
type Claims struct {
	UserID string `json:"sub"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Service provides methods for JWT token handling
type Service interface {
	GenerateToken(userID, role string) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
}

type service struct {
	jwtSecret     string
	expiryMinutes int
}

// NewService creates a new JWT service
func NewService(jwtSecret string, expiryMinutes int) Service {
	return &service{
		jwtSecret:     jwtSecret,
		expiryMinutes: expiryMinutes,
	}
}

// GenerateToken creates a new JWT token
func (s *service) GenerateToken(userID, role string) (string, error) {
	expirationTime := time.Now().Add(time.Duration(s.expiryMinutes) * time.Minute)
	
	// Create the JWT claims, which includes the user ID and expiry time
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "banking-core-mock",
		},
	}
	
	// Create the token using the claims and sign it with the secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	
	return tokenString, err
}

// ValidateToken validates the JWT token and returns the claims if valid
func (s *service) ValidateToken(tokenString string) (*Claims, error) {
	// Parse the JWT string and store the result in a claims variable
	claims := &Claims{}
	
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})
	
	if err != nil {
		return nil, err
	}
	
	// Check if the token is valid
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	
	return claims, nil
}