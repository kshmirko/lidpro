package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)


func SetupRoutes(app *fiber.App){
	app.Use(
		logger.New(),
	)
	MakePublicRoutes(app)
}