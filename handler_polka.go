package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/w0/chirpy/internal/auth"
	"github.com/w0/chirpy/internal/database"
)

func (cfg *apiConfig) handlerAddSub(w http.ResponseWriter, req *http.Request) {
	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "", err)
		return
	}

	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "invalid", err)
		return
	}

	type polkaRequest struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(req.Body)
	var polka polkaRequest
	err = decoder.Decode(&polka)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "JSON decode error", err)
		return
	}

	// we only care about the user.upgraded event.
	if polka.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}

	userID, err := uuid.Parse(polka.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid uuid", err)
		return
	}

	err = cfg.dbQueries.UpdateChirpySub(req.Context(), database.UpdateChirpySubParams{
		IsChirpyRed: true,
		ID:          userID,
	})

	if err != nil {
		respondWithError(w, http.StatusNotFound, "user not found", err)
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
