package middleware

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type Auth struct {
	secret string
}

func (a *Auth) VerifyToken(tok string) (uint, error) {
	_, err := jwt.Parse(tok, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(a.secret), nil
	})
	if err != nil {
		return 0, err
	}
	return 0, err
}
