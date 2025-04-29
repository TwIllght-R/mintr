package utils

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

func ExtractClaims(tok string, secret []byte) (map[string]interface{}, error) {
	t, err := jwt.Parse(tok, func(tkn *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil || !t.Valid {
		return nil, errors.New("invalid token")
	}
	c := t.Claims.(jwt.MapClaims)
	claims := make(map[string]interface{})
	for k, v := range c {
		claims[k] = v
	}
	return claims, nil
}
