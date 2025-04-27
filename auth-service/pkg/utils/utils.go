package utils

import (
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[n.Int64()]
	}
	return string(b), nil
}

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
