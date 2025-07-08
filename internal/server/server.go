package server

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"

	"github.com/Pratam-Kalligudda/api-gateway-go/config"
)

func SetupRoutes(app *fiber.App, cfg *config.Config, handler *ProxyHandler) {
	app.Use(logger.New())
	app.Get("/health", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
	})

	for _, service := range cfg.Services {
		app.All(service.ContextPath+"/*", handler.ForwardRequest)
	}
}
