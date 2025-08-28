package player

import (
	"github.com/gofiber/fiber/v2"

	"github.com/edcnogueira/valoguard-api/internal/models"
)

func (t *Transport) getCheatStatus(c *fiber.Ctx) error {
	name := c.Params("name")
	tag := c.Params("tag")
	region := c.Query("region")

	if name == "" || tag == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name and tag are required",
		})
	}

	req := &models.AnalyzeRequest{
		Name:   name,
		Tag:    tag,
		Region: region,
	}

	resp, err := t.service.AnalyzePlayer(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"probability": resp.Probability,
	})
}
