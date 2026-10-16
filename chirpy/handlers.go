package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/Maxy-Maxwell/go-tutorials/chirpy/internal/database"
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

func (conf *apiConfig) createChirpHandler() func(resWriter http.ResponseWriter, request *http.Request) {
	return func(resWriter http.ResponseWriter, request *http.Request) {
		type chirpRequest struct {
			Body   string `json:"body"`
			UserId string `json:"user_id"`
		}

		type chirpResponse struct {
			ID        uuid.UUID `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			Body      string    `json:"body"`
			UserID    uuid.UUID `json:"user_id"`
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

		userId, err := uuid.Parse(parsedReq.UserId)
		if err != nil {
			fmt.Println("Failed to parse userId:", err)
			writeErrorResponse(500, resWriter, "Something went wrong")
			return
		}

		chirpId := uuid.New()
		currentTime := time.Now()

		fmt.Println("UserId: ", userId)

		// Store the validated and cleaned chirp
		newChirp, err := conf.db.CreateChirp(
			request.Context(),
			database.CreateChirpParams{
				ID:        chirpId,
				CreatedAt: currentTime,
				UpdatedAt: currentTime,
				Body:      cleanedBody,
				UserID:    userId,
			})

		if err != nil {
			fmt.Println("Failed to create chirp:", err)
			writeErrorResponse(500, resWriter, "Something went wrong")
			return
		}

		writeResponseBody(201, resWriter, chirpResponse{
			ID:        newChirp.ID,
			CreatedAt: newChirp.CreatedAt,
			UpdatedAt: newChirp.UpdatedAt,
			Body:      newChirp.Body,
			UserID:    newChirp.UserID,
		})
	}
}

func (conf *apiConfig) getChirpsHandler() func(resWriter http.ResponseWriter, request *http.Request) {
	return func(resWriter http.ResponseWriter, request *http.Request) {
		type getChirpsResponse struct {
			ID        uuid.UUID `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			Body      string    `json:"body"`
			UserID    uuid.UUID `json:"user_id"`
		}

		// Query chirps
		allChirps, err := conf.db.GetChirps(request.Context())
		if err != nil {
			fmt.Println("Failed to get chirps:", err)
			writeErrorResponse(500, resWriter, "Something went wrong")
			return
		}

		var responseBody []getChirpsResponse

		for _, chirp := range allChirps {
			responseBody = append(responseBody, getChirpsResponse{
				ID:        chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body:      chirp.Body,
				UserID:    chirp.UserID,
			})
		}

		writeResponseBody(200, resWriter, responseBody)
	}
}

func (conf *apiConfig) getChirpHandler() func(resWriter http.ResponseWriter, request *http.Request) {
	return func(resWriter http.ResponseWriter, request *http.Request) {
		type getChirpsResponse struct {
			ID        uuid.UUID `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			Body      string    `json:"body"`
			UserID    uuid.UUID `json:"user_id"`
		}

		chirpId := request.PathValue("id")
		fmt.Println("Path:", chirpId)

		if chirpId == "" {
			fmt.Println("Not chirp id found:", chirpId)
			writeErrorResponse(500, resWriter, "Something went wrong")
			return
		}

		parsedChirpId, err := uuid.Parse(chirpId)
		if err != nil {
			fmt.Println("Failed to parse chirpId into UUID: ", chirpId)
			writeErrorResponse(500, resWriter, "Something went wrong")
			return
		}

		// Query chirps
		chirp, err := conf.db.GetChirpById(request.Context(), parsedChirpId)
		if err != nil {
			fmt.Println("Failed to get chirp:", err)
			// technically there could be a different sql error but this works for now
			writeErrorResponse(404, resWriter, "Not Found")
			return
		}

		writeResponseBody(200, resWriter, getChirpsResponse{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}
}

func (conf *apiConfig) statsHandler() func(resWriter http.ResponseWriter, request *http.Request) {
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
				conf.fileserverHits.Load(),
			),
		)
	}
}

func (conf *apiConfig) resetStatsHandler() func(resWriter http.ResponseWriter, request *http.Request) {
	return func(resWriter http.ResponseWriter, request *http.Request) {
		fmt.Println("Checking platform...")
		fmt.Printf("conf platform: %v\n", conf.platform)

		if conf.platform != "dev" {
			fmt.Println("Non-dev platform, returning forbidden error")
			writeErrorResponse(403, resWriter, "Forbidden")
			return
		}

		fmt.Println("Resetting stats...")
		conf.fileserverHits.Store(0)

		fmt.Println("Deleting users...")
		if err := conf.db.DeleteUsers(request.Context()); err != nil {
			fmt.Println("Failed to delete users: ", err)
			writeErrorResponse(500, resWriter, "Something went wrong")
			return
		}

		// Set headers
		resWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")

		// Set status code
		resWriter.WriteHeader(http.StatusOK)

		// Set body
		io.WriteString(resWriter, fmt.Sprintf("Hits: %v", conf.fileserverHits.Load()))
	}
}

func (conf *apiConfig) createUserHandler() func(resWriter http.ResponseWriter, request *http.Request) {
	return func(resWriter http.ResponseWriter, request *http.Request) {
		fmt.Println("In create user handler")

		type createUserRequest struct {
			Email string `json:"email"`
		}

		type createUserResponse struct {
			ID        uuid.UUID `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			Email     string    `json:"email"`
		}

		// Get Request Body
		decoder := json.NewDecoder(request.Body)
		parsedReq := createUserRequest{}

		if err := decoder.Decode(&parsedReq); err != nil {
			fmt.Println("Failed to parse request: ", err)
			writeErrorResponse(500, resWriter, "Something went wrong")
			return
		}

		if parsedReq.Email == "" {
			fmt.Println("Invalid request email: ", parsedReq.Email)
			writeErrorResponse(500, resWriter, "Something went wrong")
			return
		}

		userId := uuid.New()
		currentTime := time.Now()

		// Create User
		newUser, err := conf.db.CreateUser(request.Context(), database.CreateUserParams{
			ID:        userId,
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
			Email:     parsedReq.Email,
		})

		if err != nil {
			fmt.Println("Failed to create new user: ", err)
			writeErrorResponse(500, resWriter, "Something went wrong")
			return
		}

		fmt.Println("User: ", newUser)

		// Create and Marshal response body
		responseBody := createUserResponse{
			ID:        newUser.ID,
			CreatedAt: newUser.CreatedAt,
			UpdatedAt: newUser.UpdatedAt,
			Email:     newUser.Email,
		}

		responseBodyMarshalled, err := json.Marshal(responseBody)
		if err != nil {
			fmt.Println("Failed to marshal response: ", err)
			writeErrorResponse(500, resWriter, "Something went wrong")
			return
		}

		// Set header and status code
		resWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
		resWriter.WriteHeader(http.StatusCreated)

		// Set response body
		resWriter.Write(responseBodyMarshalled)
	}
}
