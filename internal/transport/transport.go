package transport

import "github.com/gofiber/fiber/v2"

type PublicRouter interface {
	InitPublicRoutes(app fiber.Router) error
}

type ProtectedRouter interface {
	InitProtectedRoutes(app fiber.Router) error
}

type Transport struct {
	routers []any
}

func New(routes ...any) (Transport, error) {
	return Transport{
		routers: routes,
	}, nil
}

func (t *Transport) InitRoutes(app *fiber.App) error {
	apiGroup := app.Group("/api")

	for _, router := range t.routers {
		publicRouter, ok := router.(PublicRouter)
		if !ok {
			continue
		}
		err := publicRouter.InitPublicRoutes(apiGroup)
		if err != nil {
			return err
		}
	}

	apiProtectedGroup := apiGroup.Group("/protected")
	for _, router := range t.routers {
		protectedRouter, ok := router.(ProtectedRouter)
		if !ok {
			continue
		}
		err := protectedRouter.InitProtectedRoutes(apiProtectedGroup)
		if err != nil {
			return err
		}
	}

	return nil
}
