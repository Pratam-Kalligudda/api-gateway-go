package server

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"

	"github.com/Pratam-Kalligudda/api-gateway-go/config"
	"github.com/Pratam-Kalligudda/api-gateway-go/internal/middleware"
)

func SetupRoutes(app *fiber.App, auth *middleware.Auth, cfg *config.Config, handler *ProxyHandler) {
	// This existing middleware is great for logging each incoming request
	app.Use(logger.New())

	log.Println("--- Setting up application routes ---")

	// Setup health check route first
	log.Println("Registering route: METHOD=[GET], PATH=[/health], AUTH=[false]")
	app.Get("/health", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
	})

	// Prepare auth middleware once to be reused
	authMiddleware := auth.AuthMiddleware()

	log.Println("Registering service routes from configuration...")
	for _, service := range cfg.Services {
		for _, route := range service.Routes {
			// Construct the full path for the route
			fullPath := service.ContextPath + route.Path

			// Log the details of the route being registered
			log.Printf("Registering route: METHOD=[%s], PATH=[%s], AUTH=[%t] for Service=[%s]",
				strings.Join(route.Methods, ","), fullPath, route.AuthRequired, service.Name)

			// IMPORTANT: Middleware must be passed BEFORE the final handler.
			if route.AuthRequired {
				app.Add(route.Methods, fullPath, handler.ForwardRequest, authMiddleware)
			} else {
				app.Add(route.Methods, fullPath, handler.ForwardRequest)
			}
		}
	}

	log.Println("--- Route setup complete ---")

	// The stack dump is very useful for a final verification of all registered routes.
	log.Println("--- Final Application Route Stack ---")
	stack := app.Stack()
	// Using "  " for cleaner indentation in logs
	marshStack, err := json.MarshalIndent(stack, "", "  ")
	if err != nil {
		log.Printf("Error marshaling app stack: %v", err)
	} else {
		log.Println(string(marshStack))
	}
	log.Println("------------------------------------")
}
