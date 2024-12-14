package app

import (
	"file.viktorir/internal/database/sqlite"
	"file.viktorir/internal/handler"
	"file.viktorir/internal/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func Run() {
	db, err := sqlite.Init()
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		BodyLimit:             50 * 1024 * 1024,
	})
	app.Use(cors.New(cors.Config{}))

	handler := handler.Init(db)

	router.Setup(app, *handler)

	err = app.Listen(":80")
	if err != nil {
		log.Fatal(err)
	}
}
