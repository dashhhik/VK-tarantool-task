package readhandler

import "github.com/gofiber/fiber/v2"

type ReadReq struct {
	Keys []string `json:"keys"`
}

type ReaderService interface {
	Read(keys []string) (interface{}, error)
}

type Handler struct {
	Reader ReaderService
}

func NewHandler(reader ReaderService) *Handler {
	return &Handler{Reader: reader}
}

func (h Handler) Read(c *fiber.Ctx) error {
	var req ReadReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	result, err := h.Reader.Read(req.Keys)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(200).JSON(result)

}
