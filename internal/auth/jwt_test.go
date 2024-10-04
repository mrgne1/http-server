package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

type config struct {
	UserId uuid.UUID
	TokenSecret string
	ExpiresIn time.Duration
}

func Setup(expiresIn time.Duration) config {
	return config{
		UserId: uuid.New(),
		TokenSecret: "allyourbasearebelongtous",
		ExpiresIn: expiresIn,
	}
}

func TestJWTCreation(t *testing.T) {
	cfg := Setup(time.Second)
	_, err := MakeJWT(cfg.UserId, cfg.TokenSecret, cfg.ExpiresIn)
	if err != nil {
		t.Errorf("Unable to create token: %v", err)
	}
}

func TestJWTValidation(t *testing.T) {
	cfg := Setup(time.Second)
	tokenString, err := MakeJWT(cfg.UserId, cfg.TokenSecret, cfg.ExpiresIn)
	if err != nil {
		t.Errorf("Unable to create token: %v", err)
	}

	userId, err := ValidateJWT(tokenString, cfg.TokenSecret)
	if err != nil {
		t.Errorf("Unable to validate token, %v", err)
	}
	
	if userId != cfg.UserId {
		t.Errorf("Incorrect user id: expected %v got %v", cfg.UserId, userId)
	}
}
func TestJWTSecret(t *testing.T) {
	cfg := Setup(time.Second)
	tokenString, err := MakeJWT(cfg.UserId, cfg.TokenSecret + "extra", cfg.ExpiresIn)
	if err != nil {
		t.Errorf("Unable to create token: %v", err)
	}

	_, err = ValidateJWT(tokenString, cfg.TokenSecret)
	if err == nil {
		t.Errorf("Token was incorrectly validated")
	}
}

func TestJWTExpiration(t *testing.T) {
	cfg := Setup(time.Millisecond * 100)

	tokenString, err := MakeJWT(cfg.UserId, cfg.TokenSecret, cfg.ExpiresIn)
	if err != nil {
		t.Errorf("Unable to create token: %v", err)
	}

	time.Sleep(time.Second)

	_, err = ValidateJWT(tokenString, cfg.TokenSecret)
	if err == nil {
		t.Errorf("Token didn't expire")
	}
	
}
