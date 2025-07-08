package server

import (
	"github.com/gofiber/fiber/v3"

	"github.com/Pratam-Kalligudda/api-gateway-go/internal/proxy"
)

type ProxyHandler struct {
	Service proxy.Service
}

func NewProxyHandler(svc proxy.Service) *ProxyHandler {
	return &ProxyHandler{Service: svc}
}

func (h *ProxyHandler) ForwardRequest(c fiber.Ctx) error {
	return h.Service.Forward(c)
}
