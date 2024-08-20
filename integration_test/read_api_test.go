package integration_test

import (
	"VK-test/application/reader"
	"VK-test/handler/readhandler"
	"bytes"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func (m *MockKeyValueRepo) Get(key string) (interface{}, error) {
	value, exists := m.data[key]
	if !exists {
		return nil, errors.New("key not found")
	}
	return value, nil
}

func setupReadApp(repo reader.KeyValueRepo) *fiber.App {
	logger := zap.NewNop()
	readerService := reader.NewReaderService(repo, logger)
	readHandler := readhandler.NewHandler(readerService, logger)

	app := fiber.New()
	app.Post("/api/read", readHandler.Read)

	return app
}

func TestSuccessfulRead(t *testing.T) {
	mockRepo := &MockKeyValueRepo{
		data: map[string]interface{}{
			"key1": "value1",
			"key2": "2",
			"key3": "3.14",
		},
	}
	app := setupReadApp(mockRepo)

	payload := `{"keys": ["key1", "key2", "key3"]}`
	req := httptest.NewRequest("POST", "/api/read", bytes.NewReader([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	expectedResponse := `{"data": {"key1": "value1", "key2": 2, "key3": 3.14}}`
	assert.JSONEq(t, expectedResponse, string(body))
}

func TestReadWithMissingKey(t *testing.T) {
	mockRepo := &MockKeyValueRepo{
		data: map[string]interface{}{
			"key1": "value1",
			"key3": "3.14",
		},
	}
	app := setupReadApp(mockRepo)

	payload := `{"keys": ["key1", "key2", "key3"]}`
	req := httptest.NewRequest("POST", "/api/read", bytes.NewReader([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	expectedError := `{"error": "Internal server error"}`
	assert.JSONEq(t, expectedError, string(body))
}

func TestInvalidRequestFormatRead(t *testing.T) {
	mockRepo := &MockKeyValueRepo{data: make(map[string]interface{})}
	app := setupReadApp(mockRepo)

	payload := `{"invalid": "format"}`
	req := httptest.NewRequest("POST", "/api/read", bytes.NewReader([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	expectedError := `{"error": "Invalid request format"}`
	assert.JSONEq(t, expectedError, string(body))
}

func TestReadWithUnsupportedType(t *testing.T) {
	mockRepo := &MockKeyValueRepo{
		data: map[string]interface{}{
			"key1": []string{"unsupported", "type"},
		},
	}
	app := setupReadApp(mockRepo)

	// Пример запроса с ключом, значение которого имеет неподдерживаемый тип
	payload := `{"keys": ["key1"]}`
	req := httptest.NewRequest("POST", "/api/read", bytes.NewReader([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	expectedError := `{"error": "Internal server error"}`
	assert.JSONEq(t, expectedError, string(body))
}
