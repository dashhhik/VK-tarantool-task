package integration_test

import (
	"VK-test/application/writer"
	"VK-test/handler/writehandler"
	"bytes"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net/http/httptest"
	"testing"
)

type MockKeyValueRepo struct {
	data map[string]interface{}
}

func (m *MockKeyValueRepo) Set(key string, value interface{}) error {
	if key == "errorKey" {
		return errors.New("simulated error")
	}
	m.data[key] = value
	return nil
}

func setupWriteApp(repo writer.KeyValueRepo) *fiber.App {
	logger := zap.NewNop()
	writerService := writer.NewWriterService(repo)
	writeHandler := writehandler.NewHandler(writerService, logger)

	app := fiber.New()
	app.Post("/api/write", writeHandler.Write)

	return app
}

func TestSuccessfulWrite(t *testing.T) {
	mockRepo := &MockKeyValueRepo{data: make(map[string]interface{})}
	app := setupWriteApp(mockRepo)

	payload := `{"data": {"key1": "value1", "key2": 2, "key3": 3.14}}`
	req := httptest.NewRequest("POST", "/api/write", bytes.NewReader([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := GetResponseBody(resp)
	require.NoError(t, err)

	assert.JSONEq(t, `{"status": "success"}`, body)

	assert.Equal(t, "value1", mockRepo.data["key1"])

	assert.Equal(t, float64(2), mockRepo.data["key2"])
	assert.Equal(t, 3.14, mockRepo.data["key3"])
}

func TestWriteWithError(t *testing.T) {
	mockRepo := &MockKeyValueRepo{data: make(map[string]interface{})}
	app := setupWriteApp(mockRepo)

	payload := `{"data": {"key1": "value1", "errorKey": "value2"}}`
	req := httptest.NewRequest("POST", "/api/write", bytes.NewReader([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	body, err := GetResponseBody(resp)
	require.NoError(t, err)

	assert.JSONEq(t, `{"error": "Internal server error"}`, body)

	assert.Equal(t, "value1", mockRepo.data["key1"])
	_, exists := mockRepo.data["errorKey"]
	assert.False(t, exists)
}

func TestInvalidRequestFormat(t *testing.T) {
	mockRepo := &MockKeyValueRepo{data: make(map[string]interface{})}
	app := setupWriteApp(mockRepo)

	payload := `{"data": "invalid format"}`
	req := httptest.NewRequest("POST", "/api/write", bytes.NewReader([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, err := GetResponseBody(resp)
	require.NoError(t, err)

	assert.JSONEq(t, `{"error": "Invalid 'data' field structure"}`, body)
}

func TestInvalidJSONStructure(t *testing.T) {
	mockRepo := &MockKeyValueRepo{data: make(map[string]interface{})}
	app := setupWriteApp(mockRepo)

	payload := `{"invalid JSON"`
	req := httptest.NewRequest("POST", "/api/write", bytes.NewReader([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, err := GetResponseBody(resp)
	require.NoError(t, err)

	assert.JSONEq(t, `{"error": "Invalid request format"}`, body)
}
