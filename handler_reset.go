package main

import "net/http"

func (cfg *apiConfig) handlerResetMetrics(w http.ResponseWriter, req *http.Request) {

	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "forbidden", nil)
		return
	}

	cfg.fileserverHits.Store(0)
	cfg.dbQueries.DeleteUsers(req.Context())

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	req.Write(w)

}
