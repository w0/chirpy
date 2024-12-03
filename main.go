package main

import (
	"fmt"
	"io"
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
	mux.HandleFunc("GET /api/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /api/reset", apiCfg.handlerResetMetrics)

	server := http.Server{
		Handler: mux,
		Addr:    httpPort,
	}

	server.ListenAndServe()

}

func handlerHealthz(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "OK")
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	hits := fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())
	io.WriteString(w, hits)
}

func (cfg *apiConfig) handlerResetMetrics(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	req.Write(w)

}
