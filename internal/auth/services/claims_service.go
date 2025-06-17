package services

import (
	"fmt"
	"soft-hsm/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ClaimsService struct {
	cfg   *config.Config
	Email string `json:"email"`
	Id    int    `json:"id"`
	Otp   string `json:"otp"`
	jwt.RegisteredClaims
}

func NewClaimsService(cfg *config.Config) *ClaimsService {
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
	expiresInMinutes := c.cfg.JWTConfig.ActivationExpires // Проверяем, что путь верный!
	fmt.Println("Expires: ", expiresInMinutes)
	claims := &ClaimsService{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiresInMinutes) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// Валидация активационного токена
func (c *ClaimsService) ValidateActivationToken(tokenString string) (*ClaimsService, error) {
	jwtSecret := []byte(c.cfg.JWTConfig.ActivationSecret)
	expiresInMinutes := c.cfg.JWTConfig.ActivationExpires // Проверяем, что путь верный!
	fmt.Println("Expires: ", expiresInMinutes)
	token, err := jwt.ParseWithClaims(tokenString, &ClaimsService{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		fmt.Println("TEST", err)
		return nil, fmt.Errorf("%w:", err)
	}

	claims, ok := token.Claims.(*ClaimsService)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}

func (c *ClaimsService) GenerateOTP(userID int64, code string, email string) (string, error) {
	sessionSecret := []byte(c.cfg.JWTConfig.SessionSecret)
	sessionInMinutes := c.cfg.JWTConfig.SessionExpires

	claims := &ClaimsService{
		Id:    int(userID),
		Otp:   code,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(sessionInMinutes) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(sessionSecret)
}

func (c *ClaimsService) ValidateOTPToken(tokenString string) (*ClaimsService, error) {
	jwtSecret := []byte(c.cfg.JWTConfig.SessionSecret)
	expiresInMinutes := c.cfg.JWTConfig.SessionExpires // Проверяем, что путь верный!
	fmt.Println("Expires: ", expiresInMinutes)
	token, err := jwt.ParseWithClaims(tokenString, &ClaimsService{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		fmt.Println("TEST", err)
		return nil, fmt.Errorf("%w:", err)
	}

	claims, ok := token.Claims.(*ClaimsService)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	fmt.Println(claims)
	fmt.Println(claims.Otp)

	fmt.Println(claims.Id)
	fmt.Println(claims.Email)

	return claims, nil
}
