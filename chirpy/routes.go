package main

import (
	"fmt"
	"net/http"
)

func setupRoutes(mux *http.ServeMux, stats *apiConfig) error {
	fmt.Println("Setting up file server at '/'")

	mux.Handle("/app/", stats.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir("./app")))))

	mux.HandleFunc("GET /api/healthz", readinessHandler())

	mux.HandleFunc("GET /admin/metrics", stats.statsHandler())

	mux.HandleFunc("POST /admin/reset", stats.resetStatsHandler())

	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler())

	return nil
}
