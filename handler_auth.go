package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/w0/chirpy/internal/auth"
	"github.com/w0/chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	type userLogin struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

	jwt, err := auth.MakeJWT(dbUser.ID, cfg.secret, time.Hour*1)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create jwt", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create refresh_token", err)
		return
	}

	dbRefreshToken, err := cfg.dbQueries.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    dbUser.ID,
		ExpiresAt: time.Now().AddDate(0, 0, 60),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed storing refresh token", err)
		return
	}

	type userToken struct {
		Id           uuid.UUID `json:"id"`
		Email        string    `json:"email"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
	}

	respondWithJSON(w, http.StatusOK, userToken{
		Id:           dbUser.ID,
		Email:        dbUser.Email,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
		Token:        jwt,
		RefreshToken: dbRefreshToken.Token,
	})
}

func (cfg *apiConfig) handlerRefreshJWT(w http.ResponseWriter, req *http.Request) {
	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "refresh token not found", err)
	}

	dbRefreshToken, err := cfg.dbQueries.GetRefreshToken(req.Context(), bearerToken)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "token not found", err)
		return
	}

	if dbRefreshToken.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "token expired", nil)
		return
	}

	if dbRefreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "token revoked", nil)
		return
	}

	jwt, err := auth.MakeJWT(dbRefreshToken.UserID, cfg.secret, time.Hour*1)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "create token failed", err)
		return
	}

	type refresh struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, refresh{
		Token: jwt,
	})

}

func (cfg *apiConfig) handlerRevokeRefreshToken(w http.ResponseWriter, req *http.Request) {
	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "", err)
		return
	}

	dbRefreshToken, err := cfg.dbQueries.GetRefreshToken(req.Context(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "token not found", err)
		return
	}

	now := time.Now()

	err = cfg.dbQueries.SetRevokedAt(req.Context(), database.SetRevokedAtParams{
		RevokedAt: sql.NullTime{
			Time:  now,
			Valid: true,
		},
		UpdatedAt: now,
		Token:     dbRefreshToken.Token,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to set revoked_at time", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
