package helpers

import (
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/salamanderman234/pos-backend/config"
)

func JWTParseToken(token string, secret ...string) (jwt.MapClaims, error) {
	key := strings.Join(secret, "")
	tkn, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return nil, config.ErrInvalidToken
	}
	if !tkn.Valid {
		return nil, config.ErrInvalidToken
	}
	claims := tkn.Claims.(jwt.MapClaims)
	return claims, nil
}

func JWTCreateToken(userID any, fullname string, expiryTime time.Duration, keys ...string) (string, error) {
	key := strings.Join(keys, "")
	at := time.Now()
	claims := jwt.MapClaims{
		config.AUTH_TOKEN_ID_KEY:   userID,
		config.AUTH_TOKEN_NAME_KEY: fullname,
		"iat":                      at.Unix(),
	}
	if expiryTime != 0 {
		expiredAt := at.Add(expiryTime)
		claims["exp"] = expiredAt.Unix()
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(key))
}
