package main

import (
	"encoding/json"
	"net/http"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/w0/chirpy/internal/auth"
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

	bearerToken, err := auth.GetBearerToken(req.Header)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "missing Authorization header", err)
		return
	}

	userID, err := auth.ValidateJWT(bearerToken, cfg.secret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "", err)
		return
	}

	type newChirp struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(req.Body)
	c := newChirp{}
	err = decoder.Decode(&c)

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
			UserID: userID,
		})

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to create new chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		Id:        dbChrip.ID,
		CreatedAt: dbChrip.CreatedAt,
		UpdatedAt: dbChrip.UpdatedAt,
		Body:      dbChrip.Body,
		UserId:    dbChrip.UserID,
	})

}

func cleanBody(c string) string {
	r := regexp.MustCompile("([Kk]erfuffle)|([Ss]harbert)|([Ff]ornax)")

	return r.ReplaceAllLiteralString(c, "****")
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, req *http.Request) {
	dbChirps, err := cfg.dbQueries.GetChirps(req.Context())

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed getting chirps from database", err)
		return
	}

	c := []Chirp{}

	for _, item := range dbChirps {
		c = append(c, Chirp{
			Id:        item.ID,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
			Body:      item.Body,
			UserId:    item.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, c)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, req *http.Request) {
	reqUUID, err := uuid.Parse(req.PathValue("chirpID"))

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid UUID format", err)
		return
	}

	dbChirp, err := cfg.dbQueries.GetChirp(req.Context(), reqUUID)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "failed to find chirp id", err)
		return
	}

	c := Chirp{
		Id:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserId:    dbChirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, c)
}
