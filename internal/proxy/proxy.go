package proxy

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/valyala/fasthttp"

	"github.com/Pratam-Kalligudda/api-gateway-go/config"
)

type Service interface {
	Forward(c fiber.Ctx) error
}

type proxyService struct {
	cfg        *config.Config
	httpClient *fasthttp.Client
}

func NewService(cfg *config.Config) Service {
	return &proxyService{
		cfg:        cfg,
		httpClient: &fasthttp.Client{},
	}
}

func (svc *proxyService) Forward(c fiber.Ctx) error {
	serviceConfig, err := svc.cfg.GetUrlFromEndpoint(c.Path())
	if err != nil {
		return err
	}
	downStreamReq := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(downStreamReq)

	c.Request().CopyTo(downStreamReq)

	originalPath := string(c.Request().URI().Path())
	stripedPath := strings.TrimPrefix(originalPath, serviceConfig.ContextPath)
	if !strings.HasPrefix(stripedPath, "/") {
		stripedPath = "/" + stripedPath
	}

	targetURL := fmt.Sprintf("%s%s", serviceConfig.TargetUrl, stripedPath)
	downStreamReq.SetRequestURI(targetURL)
	downStreamReq.SetHost(serviceConfig.TargetUrl)

	downStreamResp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(downStreamResp)

	if err := svc.httpClient.Do(downStreamReq, downStreamResp); err != nil {
		log.Printf("Error calling downstream service %s : %v", serviceConfig.Name, err)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": "Could not connect to the downstream service."})
	}
	c.Status(downStreamResp.StatusCode())
	downStreamResp.Header.VisitAll(func(key, value []byte) { c.Set(string(key), string(value)) })
	return c.Send(downStreamResp.Body())
}
