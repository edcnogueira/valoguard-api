package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/edcnogueira/valoguard-api/internal/providers/henrik"
	"github.com/edcnogueira/valoguard-api/internal/service/analysisservice"
	"github.com/edcnogueira/valoguard-api/internal/transport"
	"github.com/edcnogueira/valoguard-api/internal/transport/player"
)

func main() {
	henrikKey := os.Getenv("HENRIK_API_KEY")
	if henrikKey == "" {
		log.Fatal("Defina HENRIK_API_KEY como env var")
	}

	henrikClient := henrik.New(http.DefaultClient, henrikKey)

	analysisService := analysisservice.New(henrikClient)

	playerTransport := player.New(&analysisService)
	
	httpTransport, err := transport.New(&playerTransport)
	if err != nil {
		log.Fatalf("Error creating transport: %v", err)
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET",
	}))

	err = httpTransport.InitRoutes(app)
	if err != nil {
		log.Fatalf("Error initializing routes: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
