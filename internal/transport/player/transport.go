package player

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"github.com/edcnogueira/valoguard-api/internal/models"
)

type Transport struct {
	service Service
}

type Service interface {
	AnalyzePlayer(ctx context.Context, req *models.AnalyzeRequest) (*models.AnalyzeResponse, error)
}

func New(service Service) Transport {
	return Transport{
		service: service,
	}
}

func (t *Transport) InitPublicRoutes(app fiber.Router) error {
	apiGroup := app.Group("/player")

	apiGroup.Get("/cheat-status/:name/:tag", t.getCheatStatus)
	return nil
}
