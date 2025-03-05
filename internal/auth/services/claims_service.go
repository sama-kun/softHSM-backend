package services

import (
	"soft-hsm/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ClaimsService struct {
	cfg   *config.Config
	Email string `json:"email"`
	Id    int
	jwt.RegisteredClaims
}

func NewCliamsService(cfg *config.Config) *ClaimsService {
	return &ClaimsService{cfg: cfg}
}

func (c *ClaimsService) GenerateToken(id int, email string) (string, error) {
	jwtSecret := []byte(c.cfg.JWTConfig.Secret)
	expiresInMinutes := c.cfg.JWTConfig.Expires
	claims := &ClaimsService{
		Email: email,
		Id:    id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiresInMinutes) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func (s *ClaimsService) ValidateToken(tokenString string) (*ClaimsService, error) {
	token, err := jwt.ParseWithClaims(tokenString, &ClaimsService{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWTConfig.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*ClaimsService)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	return claims, nil
}

func (c *ClaimsService) GenerateActivationToken(email string) (string, error) {
	jwtSecret := []byte(c.cfg.JWTConfig.ActivationSecret)
	expiresInMinutes := c.cfg.ActivationExpires // 24 часа

	claims := jwt.MapClaims{
		"email":   email,
		"purpose": "activation",
		"exp":     time.Now().Add(time.Duration(expiresInMinutes) * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
