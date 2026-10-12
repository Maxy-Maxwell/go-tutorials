package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func readinessHandler() func(resWriter http.ResponseWriter, request *http.Request) {
	return func(resWriter http.ResponseWriter, request *http.Request) {
		fmt.Println("In readiness handler")

		// Set headers
		resWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")

		// Set status code
		resWriter.WriteHeader(http.StatusOK)

		// Set body
		io.WriteString(resWriter, "OK")
	}
}

func validateChirpHandler() func(resWriter http.ResponseWriter, request *http.Request) {
	return func(resWriter http.ResponseWriter, request *http.Request) {
		type chirpRequest struct {
			Body string `json:"body"`
		}

		type chirpResponse struct {
			CleanedBody string `json:"cleaned_body"`
		}

		decoder := json.NewDecoder(request.Body)
		parsedReq := chirpRequest{}

		if err := decoder.Decode(&parsedReq); err != nil {
			fmt.Println("Failed to parse request: ", err)
			writeErrorResponse(500, resWriter, "Something went wrong")
			return
		}

		if len(parsedReq.Body) > 140 {
			fmt.Println("Chirp too long:", parsedReq.Body)
			writeErrorResponse(400, resWriter, "Chirp is too long")
			return
		}

		cleanedBody := ""

		for _, word := range strings.Split(parsedReq.Body, " ") {
			compWord := strings.ToLower(word)
			bannedWords := []string{"kerfuffle", "sharbert", "fornax"}

			for _, bannedWord := range bannedWords {
				if compWord == bannedWord {
					word = "****"
				}
			}

			if cleanedBody == "" {
				cleanedBody = word
			} else {
				cleanedBody = cleanedBody + " " + word
			}

		}

		writeResponseBody(200, resWriter, chirpResponse{
			CleanedBody: cleanedBody,
		})
	}
}

func (stats *apiConfig) statsHandler() func(resWriter http.ResponseWriter, request *http.Request) {
	return func(resWriter http.ResponseWriter, request *http.Request) {
		fmt.Println("In stats handler")

		// Set headers
		resWriter.Header().Set("Content-Type", "text/html; charset=utf-8")

		// Set status code
		resWriter.WriteHeader(http.StatusOK)

		// Set body
		io.WriteString(
			resWriter,
			fmt.Sprintf(`
				<html>
				  <body>
				    <h1>Welcome, Chirpy Admin</h1>
				    <p>Chirpy has been visited %d times!</p>
				  </body>
				</html>`,
				stats.fileserverHits.Load(),
			),
		)
	}
}

func (stats *apiConfig) resetStatsHandler() func(resWriter http.ResponseWriter, request *http.Request) {
	return func(resWriter http.ResponseWriter, request *http.Request) {
		fmt.Println("Resetting stats...")

		stats.fileserverHits.Store(0)

		// Set headers
		resWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")

		// Set status code
		resWriter.WriteHeader(http.StatusOK)

		// Set body
		io.WriteString(resWriter, fmt.Sprintf("Hits: %v", stats.fileserverHits.Load()))
	}
}
