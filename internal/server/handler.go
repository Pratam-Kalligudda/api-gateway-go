package server

import "github.com/Pratam-Kalligudda/api-gateway-go/internal/proxy"

type APIHandler struct {
	proxy *proxy.Proxy
}

func SetupHandler(server Server) {
	app := server.App
	proxy := proxy.Proxy{}
	handler := APIHandler{proxy: &proxy}
	pubRoutes := app.Group("/api")
	userRoutes := pubRoutes.Group("/user")
	protUserRoutes := userRoutes.Group("/")
	productRoutes := pubRoutes.Group("/")
	protProdRoutes := productRoutes.Group("/")
	orderRoutes := pubRoutes.Group("/order")
	protOrderRoutes := orderRoutes.Group("/")
}
