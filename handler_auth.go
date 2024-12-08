package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/w0/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	type userLogin struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(req.Body)
	var login userLogin
	err := decoder.Decode(&login)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "JSON decode error", err)
		return
	}

	dbUser, err := cfg.dbQueries.GetUserByEmail(req.Context(), login.Email)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found", err)
	}

	err = auth.CheckPasswordHash(login.Password, dbUser.HashedPassword)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid password", err)
		return
	}

	ttl := getExpirationTime(login.ExpiresInSeconds)

	jwt, err := auth.MakeJWT(dbUser.ID, cfg.secret, ttl)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create jwt", err)
		return
	}

	type userToken struct {
		Id        uuid.UUID `json:"id"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Token     string    `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, userToken{
		Id:        dbUser.ID,
		Email:     dbUser.Email,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Token:     jwt,
	})
}

func getExpirationTime(seconds int) time.Duration {
	if seconds == 0 {
		return time.Second * 3600
	}

	if seconds > 3600 {
		return time.Second * 3600
	}

	return time.Duration(seconds)
}
