package server

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"

	"github.com/Pratam-Kalligudda/api-gateway-go/config"
	"github.com/Pratam-Kalligudda/api-gateway-go/internal/middleware"
)

func SetupRoutes(app *fiber.App, auth *middleware.Auth,cfg *config.Config, handler *ProxyHandler) {
	app.Use(logger.New())

	app.Get("/health", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
	})

	authMiddleware := auth.AuthMiddleware()

	for _, service := range cfg.Services {
		for _, route := range service.Routes {
			fullPath := service.ContextPath + route.Path
			if route.AuthRequired {
				app.Add(route.Methods,fullPath,handler.ForwardRequest,authMiddleware)
			}else{
				app.Add(route.Methods,fullPath,handler.ForwardRequest)
				
			}
		}
	}
}
