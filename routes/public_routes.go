package routes

import (
	"encoding/json"
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

		dt, _ := time.Parse("2006-01-02 15:04", datetime)

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

		// Готовим данные
		meas_avg := make([]models.Measurement, len(avg))

		for i := range meas_avg {
			dat_tmp, err := json.Marshal(avg[i].Dat)
			if err != nil {
				log.Panic(err)
			}
			dak_tmp, err := json.Marshal(avg[i].Dak)
			if err != nil {
				log.Panic(err)
			}
			meas_avg[i] = models.Measurement{
				ProfileTime:  avg[i].DateTime,
				RepRate:      uint32(avg[i].RepRate),
				ProfLen:      uint32(avg[i].ProfLen),
				ProfCnt:      uint32(avg[i].Count),
				ProfDat:      avg[i].Dat,
				ProfDak:      avg[i].Dak,
				ProfDataDat:  string(dat_tmp),
				ProfDataDak:  string(dak_tmp),
				ExperimentId: id,
			}
			models.CreateMeasurement(meas_avg[i])
		}

		// ----------------------------------------------------------------------------------------
		return c.Render("upload", fiber.Map{"status": "Загрузка данных произведена успешно!"})
	})

	publ.Get("/view", func(c *fiber.Ctx) error {
		m := models.GetAllExperimentsWithoutArchive()
		return c.Render("view", fiber.Map{"exp": m})
	})

	publ.Get("/experiment/:id", func(c *fiber.Ctx) error {
		id_s := c.Params("id")
		id, _ := strconv.ParseInt(id_s, 10, 32)

		curr_exp, err := models.GetExperimentById(id)
		if err != nil {
			log.Fatal(err)
		}
		meas, err := models.GetMeasurementsByExperimentId(id)
		if err != nil {
			log.Fatal(err)
		}

		arr := make([]models.MeasPlot, len(meas[0].ProfDak))
		for i, _ := range arr {
			arr[i].Alt = float64(i) * float64(curr_exp.VertRes)
			arr[i].Ch1 = meas[0].ProfDat[i]
			arr[i].Ch2 = meas[0].ProfDak[i]
		}

		//m := models.GetAllExperimentsWithoutArchive()
		return c.Render("detail", fiber.Map{"exp": curr_exp, "meas": meas, "arr": arr})
	})
}
