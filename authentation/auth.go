package auth

import (
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/lakshay88/reward-management-system/config"
)

var (
	cfg    *config.AppConfig
	jwtKey []byte
)

func init() {
	var err error
	cfg, err = config.LoadConfiguration("./config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
		return
	}
	jwtKey = []byte(cfg.JWTSecret)
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateTokens(email string) (string, string, error) {

	// Access Token logic
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		Username: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(cfg.AccessTokeTime) * time.Minute).Unix(),
		},
	})
	accessTokenString, err := accessToken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	// Refresh Token logic
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		Username: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(cfg.RefreshTokenTime) * 24 * time.Hour).Unix(),
		},
	})
	refreshTokenString, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, err
	}
	return claims, nil
}
