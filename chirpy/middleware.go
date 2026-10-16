package main

import (
	"net/http"
	"sync/atomic"

	"github.com/Maxy-Maxwell/go-tutorials/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

func initializeApiConfig(dbQueries *database.Queries, platform string) apiConfig {
	return apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
	}
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	// Anything that happens here, is run when the hanlder is first initialized (in setupRoutes)

	return http.HandlerFunc(func(resWriter http.ResponseWriter, request *http.Request) {
		// Anything here happens whenever the handler is called (aka calls to the /app/)
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(resWriter, request)
	})
}
