package server

import (
	"github.com/gofiber/fiber/v3"

	"github.com/Pratam-Kalligudda/api-gateway-go/config"
)

type Server struct {
	App    *fiber.App
	Config *config.Config
}
