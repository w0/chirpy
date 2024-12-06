package main

import (
	"encoding/json"
	"net/http"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/w0/chirpy/internal/database"
)

type Chirp struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerNewChirp(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	c := Chirp{}
	err := decoder.Decode(&c)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "JSON decode error", err)
		return
	}

	if len(c.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	c.Body = cleanBody(c.Body)

	dbChrip, err := cfg.dbQueries.NewChirp(req.Context(),
		database.NewChirpParams{
			Body:   c.Body,
			UserID: c.UserId,
		})

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to create new chirp", err)
		return
	}

	c.UpdatedAt = dbChrip.UpdatedAt
	c.CreatedAt = dbChrip.CreatedAt
	c.Id = dbChrip.ID

	respondWithJSON(w, http.StatusCreated, c)

}

func cleanBody(c string) string {
	r := regexp.MustCompile("([Kk]erfuffle)|([Ss]harbert)|([Ff]ornax)")

	return r.ReplaceAllLiteralString(c, "****")
}
