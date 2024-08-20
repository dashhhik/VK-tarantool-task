package readhandler

import (
	"VK-test/core"
	"errors"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type ReadReq struct {
	Keys []string `json:"keys"`
}

type ReaderService interface {
	Read(keys []string) (interface{}, error)
}

type Handler struct {
	Reader ReaderService
	Logger *zap.Logger
}

func NewHandler(reader ReaderService, logger *zap.Logger) *Handler {
	return &Handler{Reader: reader, Logger: logger}
}

func (h *Handler) Read(c *fiber.Ctx) error {
	var req ReadReq
	if err := c.BodyParser(&req); err != nil {
		h.Logger.Error("Error parsing request", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	if len(req.Keys) == 0 {
		h.Logger.Error("No keys provided in request")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	h.Logger.Info("Read request received", zap.Any("keys", req.Keys))

	result, err := h.Reader.Read(req.Keys)
	if err != nil {
		var customErr *core.CustomError
		if errors.As(err, &customErr) {
			h.Logger.Warn("Read error", zap.Error(err), zap.Any("keys", req.Keys))
			return c.Status(customErr.Code).JSON(fiber.Map{"error": customErr.Message})
		}

		h.Logger.Error("Internal error reading data", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	h.Logger.Info("Read request successful", zap.Any("result", result))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}
