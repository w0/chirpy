package main

import "net/http"

func (cfg *apiConfig) handlerResetMetrics(w http.ResponseWriter, req *http.Request) {

	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "forbidden", nil)
		return
	}

	cfg.fileserverHits.Store(0)
	err := cfg.dbQueries.DeleteUsers(req.Context())

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to reset users table", err)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	req.Write(w)

}
