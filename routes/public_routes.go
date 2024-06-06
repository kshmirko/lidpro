package routes

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/kshmirko/lidpro/utils"
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
		vert_res_s := c.FormValue("experiment-vertres", "1500.0")
		vert_res, err := strconv.ParseFloat(vert_res_s, 32)
		if err != nil {
			log.Println(err)
		}
		meas, err := utils.ReadZippedLidarArchive(mpf, vert_res)
		log.Println(len(meas))
		if err != nil {
			return c.Render("upload", fiber.Map{"status": "Какие-то проблемы с загрузкой файла!"})
		}

		// ----------------------------------------------------------------------------------------
		return c.Render("upload", fiber.Map{"status": "Загрузка данных произведена успешно!"})
	})
}
