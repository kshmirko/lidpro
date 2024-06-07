package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	app.Use(
		logger.New(),
	)
	app.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			"admin": "123123",
		},
	}))
	MakeStaticRoutes(app)
	MakePublicRoutes(app)
}
