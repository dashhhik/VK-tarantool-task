package tokens

import (
	"VK-test/core"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var JwtKey = []byte("secret_key")

type Payload struct {
	core.User
	jwt.RegisteredClaims
}

type PayloadService struct {
}

func NewPayloadService() *PayloadService {
	return &PayloadService{}
}

func (p *PayloadService) GenerateToken(userPayload core.User) (string, error) {
	claims := Payload{
		userPayload,
		jwt.RegisteredClaims{
			Subject:   "user",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (p *PayloadService) GetUserFromToken(token *jwt.Token) (*core.User, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("")
	}
	username := claims["username"].(string)
	password := claims["password"].(string)
	return &core.User{
		Username: username,
		Password: password,
	}, nil
}
