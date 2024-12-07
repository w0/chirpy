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

func (c *apiConfig) handlerNewUser(w http.ResponseWriter, req *http.Request) {

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

	dbUser, err := c.dbQueries.CreateUser(req.Context(), database.CreateUserParams{
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
