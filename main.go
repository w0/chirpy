package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/w0/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatal("Error opening database connection ", err)
	}

	httpPort := ":8080"
	serveDir := "."

	apiCfg := apiConfig{
		dbQueries: database.New(db),
		platform:  os.Getenv("PLATFORM"),
	}

	mux := http.NewServeMux()

	appHandler := http.StripPrefix("/app", http.FileServer(http.Dir(serveDir)))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(appHandler))

	mux.HandleFunc("GET /api/healthz", handlerHealthz)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerResetMetrics)
	mux.HandleFunc("POST /api/validate_chirp", apiCfg.handlerValidateChirp)
	mux.HandleFunc("POST /api/users", apiCfg.handlerNewUser)

	server := http.Server{
		Handler: mux,
		Addr:    httpPort,
	}

	server.ListenAndServe()

}
