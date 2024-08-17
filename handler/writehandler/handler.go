package writehandler

import "github.com/gofiber/fiber/v2"

type WriterService interface {
	Write(data map[string]interface{}) error
}

type Handler struct {
	Writer WriterService
}

func NewHandler(writer WriterService) *Handler {
	return &Handler{Writer: writer}
}

func (h Handler) Write(c *fiber.Ctx) error {
	var data interface{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request structure"})
	}

	innerData, ok := dataMap["data"].(map[string]interface{})
	if !ok {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid 'data' field structure"})
	}

	if err := h.Writer.Write(innerData); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(200).JSON(
		fiber.Map{"status": "success"})

}
