package middleware

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

type Auth struct {
	secret string
}

func NewAuth(secret string) *Auth {
	return &Auth{secret: secret}
}

func (a *Auth) AuthMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		headers := c.GetReqHeaders()
		auth, exist := headers["Authentication"]
		if !exist || len(auth) <= 0 || !strings.HasPrefix(auth[0], "Bearer") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "auth doesnt exists"})
		}

		token := strings.Split(auth[0], " ")
		if len(token) <= 1 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "auth doesnt exists"})
		}

		userId, err := a.VerifyToken(token[1])
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "auth doesnt exists" + err.Error()})
		}

		c.Request().Header.Add("X-User-Id", userId)

		return c.Next()
	}
}

func (a *Auth) VerifyToken(tok string) (string, error) {
	token, err := jwt.Parse(tok, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(a.secret), nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", errors.New("invalid token found")
	}

	claims := token.Claims.(jwt.MapClaims)
	userIdStr, err := claims.GetSubject()
	if err != nil {
		return "", err
	}

	return userIdStr, nil
}
