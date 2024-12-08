package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {

	p := "purified_water"

	h, err := HashPassword(p)

	if err != nil {
		t.Fatalf("failed hashing password %v", err)
	}

	err = CheckPasswordHash(p, h)

	if err != nil {
		t.Fatalf("Password doesn't match hashed value.")
	}
}

func TestNewJWT(t *testing.T) {
	userID := uuid.New()
	secret := "donthackmebro"

	_, err := MakeJWT(userID, secret, time.Minute*5)

	if err != nil {
		t.Fatalf("failed to create JWT: %v", err)
	}

}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "donthackmebro"

	jwt, err := MakeJWT(userID, secret, time.Second*5)

	if err != nil {
		t.Fatalf("failed to create JWT: %v", err)
	}

	_, err = ValidateJWT(jwt, secret)

	if err != nil {
		t.Fatalf("failed to validate jwt %v", err)
	}

}

func TestGetBearerToken(t *testing.T) {
	header := http.Header{}

	header.Add("Authorization", "Bearer 3289euihjsknlv")

	v, err := GetBearerToken(header)

	if err != nil {
		t.Fatalf("failed getting token %v", err)
	}

	if v != "3289euihjsknlv" {
		t.Fatalf("failed to trim prefix %s", v)
	}
}
