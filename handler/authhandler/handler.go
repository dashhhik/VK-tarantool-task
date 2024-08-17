package authhandler

import (
	"VK-test/core"
	"github.com/gofiber/fiber/v2"
)

type ServiceAuth interface {
	Login(user core.User) (string, error)
}

type Handler struct {
	Auth ServiceAuth
	//Logger *zap.logger
}

func NewHandler(auth ServiceAuth) *Handler {
	return &Handler{Auth: auth}

}

func (h *Handler) Login(c *fiber.Ctx) error {
	var params core.User
	if err := c.BodyParser(&params); err != nil {
		//h.Logger.Error("error parsing body: %v", err)
		return c.Status(422).JSON(fiber.Map{"errors": fiber.Map{
			"body": []string{"Invalid body"},
		}})
	}

	token, err := h.Auth.Login(params)
	if err != nil {
		//h.Logger.Error("error generating token: %v", err)
		return c.SendStatus(500)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}
