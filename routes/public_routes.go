package routes

import (
	"bytes"

	"github.com/gofiber/fiber/v2"
)

func MakePublicRoutes(app *fiber.App) {
	publ := app.Group("/public")
	publ.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})

	publ.Get("/upload", func(c *fiber.Ctx) error {
		return c.Render("upload", fiber.Map{})
	})

	publ.Post("/upload", func(c *fiber.Ctx) error {

		// Чтение файла с данными -----------------------------------------------------------------
		mpf, _ := c.FormFile("experiment-archivefile")

		rdr, err := mpf.Open()
		defer func() {
			rdr.Close()
		}()

		if err != nil {
			return c.Render("upload", fiber.Map{"status": "Какие-то проблемы с загрузкой файла!"})
		}

		buf := bytes.NewBuffer(nil)
		_, err = buf.ReadFrom(rdr)

		if err != nil {
			return c.Render("upload", fiber.Map{"status": "Какие-то проблемы при чтении файла!"})
		}
		// ----------------------------------------------------------------------------------------
		return c.Render("upload", fiber.Map{"status": "Загрузка данных произведена успешно!"})
	})
}
