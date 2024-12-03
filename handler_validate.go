package main

import (
	"encoding/json"
	"net/http"
	"regexp"
)

func (cfg *apiConfig) handlerValidateChirp(w http.ResponseWriter, req *http.Request) {

	type chirp struct {
		Body string `json:"body"`
	}

	type cleanChirp struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(req.Body)
	c := chirp{}
	err := decoder.Decode(&c)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "JSON decode error", err)
		return
	}

	if len(c.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	clean := cleanBody(c.Body)

	respondWithJSON(w, http.StatusOK, cleanChirp{CleanedBody: clean})
}

func cleanBody(c string) string {
	r := regexp.MustCompile("([Kk]erfuffle)|([Ss]harbert)|([Ff]ornax)")

	return r.ReplaceAllLiteralString(c, "****")
}
