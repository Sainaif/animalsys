package security

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sainaif/animalsys/backend/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// JWTService handles JWT token operations
type JWTService struct {
	secret               string
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

// Claims represents the JWT claims
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

// NewJWTService creates a new JWT service
func NewJWTService(secret string, accessDuration, refreshDuration time.Duration) *JWTService {
	return &JWTService{
		secret:               secret,
		accessTokenDuration:  accessDuration,
		refreshTokenDuration: refreshDuration,
	}
}

// GenerateAccessToken generates a new access token
func (s *JWTService) GenerateAccessToken(userID primitive.ObjectID, email, role string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID.Hex(),
		Email:  email,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(s.accessTokenDuration).Unix(),
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			Issuer:    "animalsys",
			Subject:   userID.Hex(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", errors.Wrap(err, 500, "failed to generate access token")
	}

	return tokenString, nil
}

// GenerateRefreshToken generates a new refresh token
func (s *JWTService) GenerateRefreshToken(userID primitive.ObjectID) (string, error) {
	now := time.Now()
	claims := jwt.StandardClaims{
		ExpiresAt: now.Add(s.refreshTokenDuration).Unix(),
		IssuedAt:  now.Unix(),
		NotBefore: now.Unix(),
		Issuer:    "animalsys",
		Subject:   userID.Hex(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", errors.Wrap(err, 500, "failed to generate refresh token")
	}

	return tokenString, nil
}

// ValidateAccessToken validates an access token and returns claims
func (s *JWTService) ValidateAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		return nil, errors.NewUnauthorized("invalid token")
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.NewUnauthorized("invalid token claims")
}

// ValidateRefreshToken validates a refresh token and returns the user ID
func (s *JWTService) ValidateRefreshToken(tokenString string) (primitive.ObjectID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		return primitive.NilObjectID, errors.NewUnauthorized("invalid refresh token")
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		userID, err := primitive.ObjectIDFromHex(claims.Subject)
		if err != nil {
			return primitive.NilObjectID, errors.NewUnauthorized("invalid user ID in token")
		}
		return userID, nil
	}

	return primitive.NilObjectID, errors.NewUnauthorized("invalid refresh token claims")
}

// ExtractUserIDFromToken extracts user ID from access token without full validation
func (s *JWTService) ExtractUserIDFromToken(tokenString string) (primitive.ObjectID, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &Claims{})
	if err != nil {
		return primitive.NilObjectID, errors.NewUnauthorized("invalid token format")
	}

	if claims, ok := token.Claims.(*Claims); ok {
		userID, err := primitive.ObjectIDFromHex(claims.UserID)
		if err != nil {
			return primitive.NilObjectID, errors.NewUnauthorized("invalid user ID in token")
		}
		return userID, nil
	}

	return primitive.NilObjectID, errors.NewUnauthorized("invalid token claims")
}
