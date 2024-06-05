package main

import (
	"embed"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/django/v3"
	"github.com/joho/godotenv"
	"github.com/kshmirko/lidpro/models"
	"github.com/kshmirko/lidpro/routes"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed templates
var viewAssets embed.FS

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки файла .env!")
	}

	meas := models.GetAllExperiments()
	log.Println(meas)

	engine := django.NewPathForwardingFileSystem(http.FS(viewAssets), "/templates", ".django")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	routes.SetupRoutes(app)
	app.Listen(os.Getenv("APP_PORT"))
}
