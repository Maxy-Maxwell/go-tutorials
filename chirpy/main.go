package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	readWriteTimeout = 5 * time.Second
)

func main() {

	if err := run(); err != nil {
		fmt.Println("Failed to execute run:", err)
		os.Exit(1)
	}
}

func run() error {

	mux := http.NewServeMux()

	fmt.Println("Setting up routes...")
	if err := setupRoutes(mux); err != nil {
		return err
	}

	//mux.Handle("/", http.FileServer(http.Dir("./served-files")))
	s := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  readWriteTimeout,
		WriteTimeout: readWriteTimeout,
	}

	fmt.Printf("Listening on http://localhost%s — serving ./served-files (GET / returns index.html)\n", s.Addr)
	return s.ListenAndServe()
}
