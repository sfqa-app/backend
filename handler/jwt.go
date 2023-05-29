package handler

import (
	"errors"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// parse jwt token and return standard claims
func ParseJwtToken(c *fiber.Ctx, token string) (*jwt.StandardClaims, error) {
	secret := os.Getenv("JWT_SECRET")

	t, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, errors.New("error parsing token")
	}

	claims, ok := t.Claims.(*jwt.StandardClaims)
	if !ok || !t.Valid || claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func GenerateToken(claims *jwt.StandardClaims) (token string, err error) {
	secret := os.Getenv("JWT_SECRET")

	c := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = c.SignedString([]byte(secret))
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}
