package auth

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var ErrNoAuthHeader error = errors.New("No Authorization Header")
var ErrInvalidAuthHeader error = errors.New("Invalid Authorization Header. Should be: 'Bearer TOKEN_STRING'")

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := jwt.StandardClaims{
		Issuer: "chirpy",
		IssuedAt: time.Now().UTC().Unix(),
		ExpiresAt: time.Now().Add(expiresIn).UTC().Unix(),
		Subject: fmt.Sprint(userID),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return  token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString string, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.StandardClaims{}
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.UUID{}, err
	}

	id := claims.Subject
	userId, err := uuid.Parse(id)
	if err != nil {
		return uuid.UUID{}, err
	}
	return userId, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	header := headers.Get("Authorization")
	if header == "" {
		return "", ErrNoAuthHeader
	}

	headerValues := strings.Split(header, " ")
	if len(headerValues) > 2 || len(headerValues) == 0 || headerValues[0] != "Bearer" {
		log.Println(header)
		return "", ErrInvalidAuthHeader
	}

	return strings.Trim(headerValues[1], " "), nil
}
