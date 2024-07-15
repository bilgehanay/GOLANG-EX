package service

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"os"
)

type UserClaims struct {
	Id       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	jwt.RegisteredClaims
}

func NewAccessToken(claims UserClaims) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return accessToken.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
}

func ParseAccessToken(accessToken string) (*UserClaims, error) {
	parsedAccessToken, err := jwt.ParseWithClaims(accessToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})

	if err != nil {
		return nil, errors.New("invalid token")
	}

	return parsedAccessToken.Claims.(*UserClaims), nil
}
