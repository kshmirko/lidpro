package routes

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kshmirko/lidpro/models"
	"github.com/kshmirko/lidpro/utils"
)

func MakePublicRoutes(app *fiber.App) {
	publ := app.Group("/public")
	publ.Get("/", func(c *fiber.Ctx) error {
		log.Println(c.AllParams())
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

		accum_time_s := c.FormValue("experiment-accumtime", "10")
		accum_time, err := strconv.Atoi(accum_time_s)

		if err != nil {
			log.Println(err)
		}

		meas, err := utils.ReadZippedLidarArchive(mpf, vert_res)
		log.Println(len(meas), err)

		avg, err := utils.MakeAverageProfiles(meas, accum_time)
		log.Println(len(avg), err)

		if err != nil {
			return c.Render("upload", fiber.Map{"status": "Какие-то проблемы с загрузкой файла!"})
		}

		datetime := c.FormValue("experiment-datetime", "2020-01-01T23:23")
		datetime = strings.ReplaceAll(datetime, "T", " ")

		log.Println(datetime)
		dt, _ := time.Parse("2006-01-02 15:04", datetime)
		log.Println(dt)

		exp := models.Experiment{
			StartTime: dt,
			Title:     c.FormValue("experiment-title", "no-title"),
			Comment:   c.FormValue("experiment-comment", "no-comment"),
			VertRes:   float32(vert_res),
			AccumTime: uint32(accum_time),
			Archive:   []byte(""),
		}
		id, err := models.CreateExperiment(exp)
		log.Println("ID=", id)
		if err != nil {
			log.Fatal(err)
		}

		// ----------------------------------------------------------------------------------------
		return c.Render("upload", fiber.Map{"status": "Загрузка данных произведена успешно!"})
	})
}
