package main

import (
	"log"

	"github.com/gofiber/fiber/v3"

	"github.com/Pratam-Kalligudda/api-gateway-go/config"
	"github.com/Pratam-Kalligudda/api-gateway-go/internal/middleware"
	"github.com/Pratam-Kalligudda/api-gateway-go/internal/proxy"
	"github.com/Pratam-Kalligudda/api-gateway-go/internal/server"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("error found : " + err.Error())
	}

	service := proxy.NewService(cfg)
	app := fiber.New()
	handler := server.NewProxyHandler(service)
	auth := middleware.NewAuth(cfg.SECRET)
	server.SetupRoutes(app,auth, cfg, handler)

	log.Println("Starting API Gateway on :8080")

	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
