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
	// Log the initial incoming request
	log.Printf("Received request to forward: %s", c.Path())

	serviceConfig, err := svc.cfg.GetUrlFromEndpoint(c.Path())
	if err != nil {
		// Log the error when no service configuration is found
		log.Printf("Error: Could not find service configuration for path %s: %v", c.Path(), err)
		return err // Or a specific fiber.Error
	}

	// Log which service the request is being routed to
	log.Printf("Routing to service: '%s' with target URL: %s", serviceConfig.Name, serviceConfig.TargetUrl)

	downStreamReq := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(downStreamReq)

	downStreamReq.SetBody(c.Body())
	downStreamReq.Header.SetMethod(c.Method())
	c.Request().Header.VisitAll(func(key, value []byte) { downStreamReq.Header.Add(string(key), string(value)) })

	userIdHeader := string(c.Request().Header.Peek("X-User-Id"))
	emailHeader := string(c.Request().Header.Peek("X-Email"))
	roleHeader := string(c.Request().Header.Peek("X-Role"))

	log.Printf("userId : [%s] | email : [%s] | role : [%s]", userIdHeader, emailHeader, roleHeader)

	downStreamReq.Header.Add("X-User-Id", userIdHeader)
	downStreamReq.Header.Add("X-Email", emailHeader)
	downStreamReq.Header.Add("X-Role", roleHeader)

	originalPath := string(c.Request().URI().Path())
	if !strings.HasPrefix(originalPath, "/") {
		originalPath = "/" + originalPath
	}

	targetURL := fmt.Sprintf("%s%s", serviceConfig.TargetUrl, originalPath)
	downStreamReq.SetRequestURI(targetURL)

	// In modern fasthttp, setting the Host header is often unnecessary as it's derived
	// from the request URI, but we'll keep it for clarity.
	// downStreamReq.SetHost(serviceConfig.TargetUrl)

	// Log the final downstream URL being called
	log.Printf("Forwarding request to downstream URL: %s", targetURL)

	downStreamResp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(downStreamResp)

	if err := svc.httpClient.Do(downStreamReq, downStreamResp); err != nil {
		// This existing log is good and specific
		log.Printf("Error calling downstream service '%s': %v", serviceConfig.Name, err)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": "Could not connect to the downstream service."})
	}

	// Log the response status from the downstream service
	log.Printf("Received response from '%s' with status: %d", serviceConfig.Name, downStreamResp.StatusCode())

	c.Status(downStreamResp.StatusCode())
	downStreamResp.Header.VisitAll(func(key, value []byte) { c.Set(string(key), string(value)) })
	return c.Send(downStreamResp.Body())
}
