package tokens

import (
	"VK-test/core"
	"fmt"
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
	// Создание payload с пользователем и стандартными claims
	claims := Payload{
		userPayload,
		jwt.RegisteredClaims{
			Subject:   "user",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Создание токена с использованием HS256 алгоритма и подписывание ключом
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		// Возвращаем обернутую ошибку с контекстом
		return "", fmt.Errorf("failed to sign JWT token: %w", err)
	}

	return tokenString, nil
}

func (p *PayloadService) GetUserFromToken(token *jwt.Token) (*core.User, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, core.NewCustomError(400, "invalid token claims")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return nil, core.NewCustomError(400, "username claim is missing or not a string")
	}

	password, ok := claims["password"].(string)
	if !ok {
		return nil, core.NewCustomError(400, "password claim is missing or not a string")
	}

	return &core.User{
		Username: username,
		Password: password,
	}, nil
}
