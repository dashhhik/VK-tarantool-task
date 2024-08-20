package writehandler

import (
	"VK-test/core"
	"errors"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type WriterService interface {
	Write(data map[string]interface{}) error
}

type Handler struct {
	Writer WriterService
	Logger *zap.Logger
}

func NewHandler(writer WriterService, logger *zap.Logger) *Handler {
	return &Handler{
		Writer: writer,
		Logger: logger,
	}
}

func (h *Handler) Write(c *fiber.Ctx) error {
	var data interface{}
	if err := c.BodyParser(&data); err != nil {
		h.Logger.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	dataMap, ok := data.(map[string]interface{})
	if !ok {
		h.Logger.Error("Request body is not a valid JSON object")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request structure"})
	}

	innerData, ok := dataMap["data"].(map[string]interface{})
	if !ok {
		h.Logger.Error("'data' field is missing or not a valid JSON object")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid 'data' field structure"})
	}

	if err := h.Writer.Write(innerData); err != nil {
		var customErr *core.CustomError
		if errors.As(err, &customErr) {
			h.Logger.Warn("Error writing data", zap.Error(err))
			return c.Status(customErr.Code).JSON(fiber.Map{"error": customErr.Message})
		}

		h.Logger.Error("Failed to write data", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	h.Logger.Info("Data written successfully")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}
