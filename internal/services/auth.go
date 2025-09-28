package services

import (
	"time"

	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	tokenAuth *jwtauth.JWTAuth
}

func NewAuthService(jwtSecret string) *AuthService {
	tokenAuth := jwtauth.New("HS256", []byte(jwtSecret), nil)

	return &AuthService{
		tokenAuth: tokenAuth,
	}
}

// GetTokenAuth returns the JWT auth instance for middleware
func (as *AuthService) GetTokenAuth() *jwtauth.JWTAuth {
	return as.tokenAuth
}

// GenerateToken creates a JWT token for a user
func (as *AuthService) GenerateToken(userID int, email, username, role string) (string, error) {
	claims := map[string]interface{}{
		"user_id":  userID,
		"email":    email,
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(), // Token expires in 24 hours
		"iat":      time.Now().Unix(),
		"iss":      "kafka-blog-api",
	}

	_, tokenString, err := as.tokenAuth.Encode(claims)
	return tokenString, err
}

// HashPassword hashes a password using bcrypt
func (as *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword compares a password with a hash
func (as *AuthService) CheckPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
