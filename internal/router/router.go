package router

import (
	"file.viktorir/internal/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"os"
)

func Setup(app *fiber.App, handler handler.FileHandler) {
	api := app.Group("/", logger.New(logger.Config{
		Format:     "{\"time\":\"${time}\",\"method\":\"${method}\",\"path\":\"${path}\",\"status\":${status},\"user-agent\":\"${ua}\",\"route\":\"${route}\",\"error\":\"${error}\"}\n",
		TimeFormat: "02/Jan/2006:15:04:05",
		TimeZone:   "Local",
		Output:     os.Stdout,
	}))

	api.Post("/upload", handler.Upload)
	api.Get("/:short_link", handler.GetByShort)
	api.Get("/file/:user_id/:type/:subtype/:name", handler.GetByFull)
}
