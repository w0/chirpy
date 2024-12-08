package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	issueTime := jwt.NewNumericDate(time.Now())
	expireTime := jwt.NewNumericDate(issueTime.Add(expiresIn))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  issueTime,
		ExpiresAt: expireTime,
		Subject:   userID.String(),
	})

	signed, err := token.SignedString([]byte(tokenSecret))

	if err != nil {
		return "", err
	}

	return signed, nil

}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok {
		return uuid.MustParse(claims.Subject), nil
	}

	return uuid.Nil, fmt.Errorf("unknown claims type")
}

func GetBearerToken(headers http.Header) (string, error) {
	if value := headers.Get("Authorization"); value != "" {
		return strings.TrimPrefix(value, "Bearer "), nil
	}

	return "", fmt.Errorf("authorization not found in headers")
}

func MakeRefreshToken() (string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}
