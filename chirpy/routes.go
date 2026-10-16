package main

import (
	"fmt"
	"net/http"
)

func setupRoutes(mux *http.ServeMux, conf *apiConfig) error {
	fmt.Println("Setting up file server at '/'")

	mux.Handle("/app/", conf.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir("./app")))))

	mux.HandleFunc("POST /api/users", conf.createUserHandler())

	mux.HandleFunc("POST /api/chirps", conf.createChirpHandler())

	mux.HandleFunc("GET /api/healthz", readinessHandler())

	mux.HandleFunc("GET /admin/metrics", conf.statsHandler())

	mux.HandleFunc("POST /admin/reset", conf.resetStatsHandler())

	return nil
}
