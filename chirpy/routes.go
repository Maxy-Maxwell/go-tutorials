package main

import (
	"fmt"
	"net/http"
)

func setupRoutes(mux *http.ServeMux, conf *apiConfig) error {
	fmt.Println("Setting up file server at '/'")

	mux.Handle("/app/", conf.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir("./app")))))

	// Users handlers
	mux.HandleFunc("POST /api/users", conf.createUserHandler())
	mux.HandleFunc("PUT /api/users", conf.middlewareJWTAuth(createUpdateUserHandler))
	mux.HandleFunc("POST /api/polka/webhooks", conf.createUpdateUserRedHandler())

	// User Login
	mux.HandleFunc("POST /api/login", conf.createLoginHandler())

	// Token handlers
	mux.HandleFunc("POST /api/refresh", conf.createRefreshTokenHandler())
	mux.HandleFunc("POST /api/revoke", conf.createRevokeTokenHandler())

	// Chirps handlers
	mux.HandleFunc("POST /api/chirps", conf.middlewareJWTAuth(createChirpHandler))
	mux.HandleFunc("GET /api/chirps", conf.getChirpsHandler())
	mux.HandleFunc("GET /api/chirps/{id}", conf.getChirpHandler())
	mux.HandleFunc("DELETE /api/chirps/{id}", conf.middlewareJWTAuth(deleteChirpHandler))

	mux.HandleFunc("GET /api/healthz", readinessHandler())

	mux.HandleFunc("GET /admin/metrics", conf.statsHandler())

	mux.HandleFunc("POST /admin/reset", conf.resetStatsHandler())

	return nil
}
