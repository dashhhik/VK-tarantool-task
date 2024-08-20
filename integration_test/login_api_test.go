package integration_test

import (
	"VK-test/application/auth"
	"VK-test/application/tokens"
	"VK-test/core"
	"VK-test/handler/authhandler"
	logger2 "VK-test/infrastructure/logger"
	"VK-test/infrastructure/repo/user_data"
	"bytes"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func GetResponseBody(resp *http.Response) (string, error) {
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bodyBytes), nil
}

func setupApp(mockTokenService *tokens.MockPayloadService, mockUserRepo *user_data.MockUserRepo) *fiber.App {
	logger := logger2.NewLogger()
	serviceAuth := auth.NewService(mockTokenService, mockUserRepo)
	authHandler := authhandler.NewHandler(serviceAuth, logger)

	app := fiber.New()
	app.Post("/api/login", authHandler.Login)

	return app
}

func TestSuccessfulLogin(t *testing.T) {
	// Настройка моков
	mockTokenService := new(tokens.MockPayloadService)
	mockUserRepo := new(user_data.MockUserRepo)

	user := core.User{
		Username: "admin",
		Password: "presale",
	}

	// Настраиваем ожидания моков
	mockUserRepo.On("Get", user.Username).Return(user.Password, nil)
	mockTokenService.On("GenerateToken", user).Return("mocked-jwt-token", nil)

	// Создаем приложение с моками
	app := setupApp(mockTokenService, mockUserRepo)

	// Создаем и отправляем запрос
	req := httptest.NewRequest("POST", "/api/login", bytes.NewReader([]byte(`{"username": "admin", "password": "presale"}`)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	// Проверяем результат
	assert.Equal(t, 200, resp.StatusCode)

	body, err := GetResponseBody(resp)
	require.NoError(t, err)

	assert.JSONEq(t, `{"token": "mocked-jwt-token"}`, body)

	// Проверяем, что моки вызваны с правильными параметрами
	mockUserRepo.AssertExpectations(t)
	mockTokenService.AssertExpectations(t)
}

func TestInvalidPassword(t *testing.T) {
	// Настройка моков
	mockTokenService := new(tokens.MockPayloadService)
	mockUserRepo := new(user_data.MockUserRepo)

	user := core.User{
		Username: "admin",
		Password: "wrongpassword",
	}

	// Настраиваем ожидания моков
	mockUserRepo.On("Get", user.Username).Return("correctpassword", nil)

	// Создаем приложение с моками
	app := setupApp(mockTokenService, mockUserRepo)

	// Создаем и отправляем запрос
	req := httptest.NewRequest("POST", "/api/login", bytes.NewReader([]byte(`{"username": "admin", "password": "wrongpassword"}`)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	// Проверяем результат
	assert.Equal(t, 401, resp.StatusCode)

	body, err := GetResponseBody(resp)
	require.NoError(t, err)

	assert.JSONEq(t, `{"error": "invalid credentials"}`, body)

	// Проверяем, что моки вызваны с правильными параметрами
	mockUserRepo.AssertExpectations(t)
	mockTokenService.AssertExpectations(t)
}

func TestWrongPayloadFormat(t *testing.T) {
	// Настройка моков
	mockTokenService := new(tokens.MockPayloadService)
	mockUserRepo := new(user_data.MockUserRepo)

	// Создаем приложение с моками
	app := setupApp(mockTokenService, mockUserRepo)

	// Создаем и отправляем запрос с неправильным форматом JSON
	req := httptest.NewRequest("POST", "/api/login", bytes.NewReader([]byte(`{"f": "5tf", "f54": "4treg"}`))) // некорректный формат JSON
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	// Проверяем результат
	assert.Equal(t, 400, resp.StatusCode)

	body, err := GetResponseBody(resp)
	require.NoError(t, err)

	assert.JSONEq(t, `{"error": "invalid request format"}`, body)

	// Проверяем, что моки вызваны с правильными параметрами
	mockUserRepo.AssertExpectations(t)
	mockTokenService.AssertExpectations(t)
}
