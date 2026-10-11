package main

import (
	"fmt"
	"net/http"
)

func setupRoutes(mux *http.ServeMux) error {
	fmt.Println("Setting up file server at '/'")

	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("./app"))))

	mux.HandleFunc("/healthz", readinessHandler())

	return nil
}
