package authhandler

import (
	"VK-test/core"
	"errors"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type ServiceAuth interface {
	Login(user core.User) (string, error)
}

type Handler struct {
	Auth   ServiceAuth
	Logger *zap.Logger
}

func NewHandler(auth ServiceAuth, logger *zap.Logger) *Handler {
	return &Handler{Auth: auth, Logger: logger}
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var params core.User
	if err := c.BodyParser(&params); err != nil {
		h.Logger.Error("failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	h.Logger.Info("login request", zap.String("username", params.Username))

	token, err := h.Auth.Login(params)
	if err != nil {
		var customErr *core.CustomError
		if ok := errors.As(err, &customErr); ok {
			h.Logger.Warn("authentication failed", zap.String("username", params.Username), zap.Error(err))
			return c.Status(customErr.Code).JSON(fiber.Map{"error": customErr.Message})
		}

		h.Logger.Error("error during login", zap.String("username", params.Username), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}
