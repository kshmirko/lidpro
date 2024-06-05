package routes

import "github.com/gofiber/fiber/v2"

func MakePublicRoutes(app *fiber.App) {
	publ := app.Group("/public")
	publ.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Fiber": "Hello World!\n Pram-pam-pam",
		})
	})
}
