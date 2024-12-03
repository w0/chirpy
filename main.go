package main

import (
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	httpPort := ":8080"
	serveDir := "."

	apiCfg := apiConfig{}

	mux := http.NewServeMux()

	appHandler := http.StripPrefix("/app", http.FileServer(http.Dir(serveDir)))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(appHandler))

	mux.HandleFunc("GET /api/healthz", handlerHealthz)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerResetMetrics)
	mux.HandleFunc("POST /api/validate_chirp", apiCfg.handlerValidateChirp)

	server := http.Server{
		Handler: mux,
		Addr:    httpPort,
	}

	server.ListenAndServe()

}
