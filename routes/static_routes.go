package routes

import "github.com/gofiber/fiber/v2"

func MakeStaticRoutes(app *fiber.App) {
	app.Static("/static", "static/")
}
