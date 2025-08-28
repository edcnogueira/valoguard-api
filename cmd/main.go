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
	// Get API key from environment
	henrikKey := os.Getenv("HENRIK_API_KEY")
	if henrikKey == "" {
		log.Fatal("Defina HENRIK_API_KEY como env var")
	}

	// Initialize providers
	henrikClient := henrik.New(http.DefaultClient, henrikKey)

	// Initialize services
	analysisService := analysisservice.New(henrikClient)

	// Initialize transport handlers
	playerTransport := player.New(&analysisService)
	
	// Create main transport with all routers
	httpTransport, err := transport.New(&playerTransport)
	if err != nil {
		log.Fatalf("Error creating transport: %v", err)
	}

	// Setup Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Setup CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET",
	}))

	// Setup routes
	err = httpTransport.InitRoutes(app)
	if err != nil {
		log.Fatalf("Error initializing routes: %v", err)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
