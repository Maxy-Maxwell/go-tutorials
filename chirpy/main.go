package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Maxy-Maxwell/go-tutorials/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	readWriteTimeout = 5 * time.Second
)

var conf apiConfig

func main() {

	// Load up the .env file, instead of just relying on process env variables
	if err := godotenv.Load(); err != nil {
		fmt.Println("load .env: ", err)
		os.Exit(1)
	}

	if err := run(); err != nil {
		fmt.Println("Failed to execute run:", err)
		os.Exit(1)
	}
}

func run() error {
	// Store paltform from .env
	platform := os.Getenv("PLATFORM")

	// Setup database connection:
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return fmt.Errorf("Failed to open DB connection: %v", err)
	}

	dbQueries := database.New(db)

	// Initialize config struct. Used for keep track of total requests received, and pass the database
	conf = initializeApiConfig(dbQueries, platform)

	mux := http.NewServeMux()

	fmt.Println("Setting up routes...")
	if err := setupRoutes(mux, &conf); err != nil {
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
