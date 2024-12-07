package main

import (
	"encoding/json"
	"net/http"

	"github.com/w0/chirpy/internal/auth"
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

	respondWithJSON(w, http.StatusOK, User{
		Id:        dbUser.ID,
		Email:     dbUser.Email,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	})
}
