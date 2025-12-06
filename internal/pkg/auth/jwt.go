package auth

import (
	"errors"
	"time"

	"TulaHackDays-Backend/internal/config"
	"TulaHackDays-Backend/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	jwt.RegisteredClaims
}

type AuthService struct {
	config *config.JWTConfig
}

func NewAuthService(cfg *config.JWTConfig) *AuthService {
	return &AuthService{config: cfg}
}

func (a *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (a *AuthService) CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (a *AuthService) GenerateToken(user *domain.User) (string, error) {
	expirationTime := time.Now().Add(time.Duration(a.config.Expiration) * time.Hour)

	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "TulaHackDays-Backend",
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(a.config.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.config.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func (a *AuthService) ExtractUserID(tokenString string) (uuid.UUID, error) {
	claims, err := a.ValidateToken(tokenString)
	if err != nil {
		return uuid.Nil, err
	}
	return claims.UserID, nil
}
