package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/w0/chirpy/internal/auth"
	"github.com/w0/chirpy/internal/database"
)

type User struct {
	Id        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (cfg *apiConfig) handlerNewUser(w http.ResponseWriter, req *http.Request) {

	type newUser struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)
	var u newUser
	err := decoder.Decode(&u)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed decoding user payload", err)
		return
	}

	hashed, err := auth.HashPassword(u.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed hashing password", err)
		return
	}

	dbUser, err := cfg.dbQueries.CreateUser(req.Context(), database.CreateUserParams{
		Email:          u.Email,
		HashedPassword: hashed,
	})

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed creating user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		Id:        dbUser.ID,
		Email:     dbUser.Email,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	})
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, req *http.Request) {
	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "authorization not found", nil)
		return
	}

	userID, err := auth.ValidateJWT(bearerToken, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "token not valid", err)
		return
	}

	dbUser, err := cfg.dbQueries.GetUserByID(req.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "user not found", err)
		return
	}

	type userUpdate struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)
	update := userUpdate{}
	err = decoder.Decode(&update)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "JSON decode error", err)
	}

	hashed, err := auth.HashPassword(update.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "", err)
		return
	}

	dbUpdate, err := cfg.dbQueries.UpdateUser(req.Context(), database.UpdateUserParams{
		Email:          update.Email,
		HashedPassword: hashed,
		ID:             dbUser.ID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to update user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{Id: dbUpdate.ID, CreatedAt: dbUpdate.CreatedAt, UpdatedAt: dbUpdate.UpdatedAt, Email: dbUpdate.Email})
}
