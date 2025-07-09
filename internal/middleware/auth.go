package middleware

import (
	"errors"
	"fmt"
	"log"
	"strconv"
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
		log.Printf("Running auth middleware for path: %s", c.Path())
		headers := c.GetReqHeaders()
		auth, exist := headers["Authorization"]
		if !exist || len(auth) <= 0 || !strings.HasPrefix(auth[0], "Bearer ") {
			// Log the specific reason for failure
			log.Println("Auth failed: 'Authentication' header missing, empty, or not a Bearer token.")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or malformed authorization header"})
		}

		tokenParts := strings.Split(auth[0], " ")
		if len(tokenParts) != 2 {
			log.Println("Auth failed: Token format is incorrect. Should be 'Bearer <token>'.")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token format"})
		}
		tokenString := tokenParts[1]

		// Log before attempting to verify
		log.Println("Attempting to verify token...")
		userId, email, role, err := a.VerifyToken(tokenString)
		if err != nil {
			// Log the actual verification error
			log.Printf("Auth failed: Token verification error: %v", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		// Log success and the authenticated user
		log.Printf("Authentication successful for user ID: %s", userId)
		c.Request().Header.Add("X-User-Id", userId)
		c.Request().Header.Add("X-Email", email)
		c.Request().Header.Add("X-Role", role)

		return c.Next()
	}
}

func (a *Auth) VerifyToken(tok string) (string, string, string, error) {
	token, err := jwt.Parse(tok, func(t *jwt.Token) (any, error) {
		// Verify the signing method is what you expect
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("Token verification error: unexpected signing method: %v", t.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(a.secret), nil
	})
	// This error is returned if the token is malformed or the signing key callback fails
	if err != nil {
		log.Printf("Token parsing failed: %v", err)
		return "", "", "", err
	}

	// Check if the token is valid after parsing
	if !token.Valid {
		log.Println("Token verification error: token is invalid.")
		return "", "", "", errors.New("invalid token")
	}

	// Log before extracting claims
	log.Println("Token is valid, extracting claims...")
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("Token verification error: could not cast claims to jwt.MapClaims.")
		return "", "", "", errors.New("error reading token claims")
	}

	userId, ok := claims["sub"].(float64)
	email := claims["email"].(string)
	role := claims["role"].(string)

	if !ok || len(email) <= 0 || len(role) <= 0 {
		log.Printf("Token verification error: could not get 'subject' or 'email' or 'role' claim: %v", err)
		return "", "", "", errors.New("could not get user ID from token : ")
	}

	userIdStr := strconv.FormatFloat(userId, 'f', -1, 64)
	log.Printf("Successfully verified token for subject (user ID): %s", userIdStr)

	return userIdStr, email, role, nil
}
