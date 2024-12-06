package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}

func (c *apiConfig) handlerNewUser(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	var u User
	err := decoder.Decode(&u)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed decoding user payload", err)
		return
	}

	dbUser, err := c.dbQueries.CreateUser(req.Context(), u.Email)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed creating user", err)
		return
	}

	u.Created_at = dbUser.CreatedAt
	u.Updated_at = dbUser.UpdatedAt
	u.Id = dbUser.ID

	respondWithJSON(w, http.StatusCreated, u)
}
