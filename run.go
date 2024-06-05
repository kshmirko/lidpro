package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/kshmirko/lidpro/models"
	"github.com/kshmirko/lidpro/routes"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки файла .env!")
	}

	meas := models.GetAllExperiments()
	log.Println(meas)

	//Setup routes facilities
	app := fiber.New()
	routes.SetupRoutes(app)
	app.Listen(os.Getenv("APP_PORT"))
}
