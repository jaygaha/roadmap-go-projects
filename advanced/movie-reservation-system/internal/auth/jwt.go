package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret     string
	ExpireTime time.Duration
}

// Claims represents the JWT claims
type Claims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// AuthService handles JWT operations
type AuthService struct {
	config *JWTConfig
}

// NewAuthService creates a new JWT service
func NewAuthService(secret string, expireHours int) *AuthService {
	return &AuthService{
		config: &JWTConfig{
			Secret:     secret,
			ExpireTime: time.Duration(expireHours) * time.Hour,
		},
	}
}

// GenerateToken generates a JWT token for the given user
func (s *AuthService) GenerateToken(userID int64, email string, role string) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.config.ExpireTime)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "movie-reservation-system",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.config.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

// GetTokenFromHeader extracts the token from the Authorization header
func (s *AuthService) GetTokenFromHeader(authHeader string) string {
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return ""
	}
	return authHeader[7:]
}
