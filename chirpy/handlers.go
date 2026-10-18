package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/Maxy-Maxwell/go-tutorials/chirpy/internal/auth"
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

func createChirpHandler(conf *apiConfig, headerUserId uuid.UUID, resWriter http.ResponseWriter, request *http.Request) {
	type chirpRequest struct {
		Body string `json:"body"`
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
		writeErrorResponse(http.StatusInternalServerError, resWriter, "Something went wrong")
		return
	}

	if len(parsedReq.Body) > 140 {
		fmt.Println("Chirp too long:", parsedReq.Body)
		writeErrorResponse(http.StatusBadRequest, resWriter, "Chirp is too long")
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

	chirpId := uuid.New()
	currentTime := time.Now()

	// Store the validated and cleaned chirp
	newChirp, err := conf.db.CreateChirp(
		request.Context(),
		database.CreateChirpParams{
			ID:        chirpId,
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
			Body:      cleanedBody,
			UserID:    headerUserId,
		})

	if err != nil {
		fmt.Println("Failed to create chirp:", err)
		writeErrorResponse(http.StatusInternalServerError, resWriter, "Something went wrong")
		return
	}

	writeResponseBody(http.StatusCreated, resWriter, chirpResponse{
		ID:        newChirp.ID,
		CreatedAt: newChirp.CreatedAt,
		UpdatedAt: newChirp.UpdatedAt,
		Body:      newChirp.Body,
		UserID:    newChirp.UserID,
	})
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
			writeErrorResponse(http.StatusInternalServerError, resWriter, "Something went wrong")
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

		writeResponseBody(http.StatusOK, resWriter, responseBody)
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

		if chirpId == "" {
			fmt.Println("Not chirp id found:", chirpId)
			writeErrorResponse(http.StatusInternalServerError, resWriter, "Something went wrong")
			return
		}

		parsedChirpId, err := uuid.Parse(chirpId)
		if err != nil {
			fmt.Println("Failed to parse chirpId into UUID: ", chirpId)
			writeErrorResponse(http.StatusInternalServerError, resWriter, "Something went wrong")
			return
		}

		// Query chirps
		chirp, err := conf.db.GetChirpById(request.Context(), parsedChirpId)
		if err != nil {
			fmt.Println("Failed to get chirp:", err)
			// technically there could be a different sql error but this works for now
			writeErrorResponse(http.StatusNotFound, resWriter, "Not Found")
			return
		}

		writeResponseBody(http.StatusOK, resWriter, getChirpsResponse{
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
		if conf.platform != "dev" {
			fmt.Println("Non-dev platform, returning forbidden error")
			writeErrorResponse(http.StatusForbidden, resWriter, "Forbidden")
			return
		}

		conf.fileserverHits.Store(0)

		if err := conf.db.DeleteUsers(request.Context()); err != nil {
			fmt.Println("Failed to delete users: ", err)
			writeErrorResponse(http.StatusInternalServerError, resWriter, "Something went wrong")
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
		type createUserRequest struct {
			Email    string `json:"email"`
			Password string `json:"password"`
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
			writeErrorResponse(http.StatusInternalServerError, resWriter, "Something went wrong")
			return
		}

		if parsedReq.Email == "" {
			fmt.Println("Invalid request email: ", parsedReq.Email)
			writeErrorResponse(http.StatusInternalServerError, resWriter, "Something went wrong")
			return
		}

		if parsedReq.Password == "" {
			fmt.Println("Invalid request password: ", parsedReq.Password)
			writeErrorResponse(http.StatusInternalServerError, resWriter, "Something went wrong")
			return
		}

		userId := uuid.New()
		currentTime := time.Now()
		hashedPassword, err := auth.HashPassword(parsedReq.Password)
		if err != nil {
			writeErrorResponse(http.StatusInternalServerError, resWriter, "Something went wrong")
			return
		}

		// Create User
		newUser, err := conf.db.CreateUser(request.Context(), database.CreateUserParams{
			ID:             userId,
			CreatedAt:      currentTime,
			UpdatedAt:      currentTime,
			Email:          parsedReq.Email,
			HashedPassword: hashedPassword,
		})

		if err != nil {
			fmt.Println("Failed to create new user: ", err)
			writeErrorResponse(http.StatusInternalServerError, resWriter, "Something went wrong")
			return
		}

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
			writeErrorResponse(http.StatusInternalServerError, resWriter, "Something went wrong")
			return
		}

		// Set header and status code
		resWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
		resWriter.WriteHeader(http.StatusCreated)

		// Set response body
		resWriter.Write(responseBodyMarshalled)
	}
}

func (conf *apiConfig) createLoginHandler() func(resWriter http.ResponseWriter, request *http.Request) {
	return func(resWriter http.ResponseWriter, request *http.Request) {
		type createUserRequest struct {
			Email                 string `json:"email"`
			Password              string `json:"password"`
			TokenExpiresInSeconds string `json:"expires_in_seconds"`
		}

		type createUserResponse struct {
			ID        uuid.UUID `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			Email     string    `json:"email"`
			Token     string    `json:"token"`
		}

		// Get Request Body
		decoder := json.NewDecoder(request.Body)
		parsedReq := createUserRequest{}

		if err := decoder.Decode(&parsedReq); err != nil {
			fmt.Println("Failed to parse request: ", err)
			writeErrorResponse(http.StatusInternalServerError, resWriter, "Something went wrong")
			return
		}

		if parsedReq.Email == "" {
			fmt.Println("Invalid request email: ", parsedReq.Email)
			writeErrorResponse(http.StatusInternalServerError, resWriter, "Something went wrong")
			return
		}

		if parsedReq.Password == "" {
			fmt.Println("Invalid request password: ", parsedReq.Password)
			writeErrorResponse(http.StatusInternalServerError, resWriter, "Something went wrong")
			return
		}

		// Grab tokenExpiresIn and default to 1 hour if > 1 hour, not provided or negative
		tokenExpiresIn, err := strconv.Atoi(parsedReq.TokenExpiresInSeconds)
		if err != nil || tokenExpiresIn > int(time.Hour.Seconds()) || tokenExpiresIn < 1 {
			tokenExpiresIn = int(time.Hour.Seconds())
		}

		// Get User by Email
		foundUser, err := conf.db.GetUserByEmail(request.Context(), parsedReq.Email)
		if err != nil {
			fmt.Println("Failed to get user by email: ", parsedReq.Email, err)
			writeErrorResponse(http.StatusInternalServerError, resWriter, "Something went wrong")
			return
		}

		// Ensure password matches hash
		passwordMatches, err := auth.CheckPasswordHash(parsedReq.Password, foundUser.HashedPassword)
		if err != nil {
			writeErrorResponse(http.StatusInternalServerError, resWriter, "Something went wrong")
			return
		}

		if !passwordMatches {
			writeErrorResponse(http.StatusUnauthorized, resWriter, "Unauthorized")
			return
		}

		// Grab a token now that we know the user is authenticated (can happen earlier, but will save work if we only do it after auth)
		jwtToken, err := auth.MakeJWT(foundUser.ID, conf.jwt_secret, time.Duration(tokenExpiresIn)*time.Second)
		if err != nil {
			fmt.Println("Failed to generate token: ", err, ". UserId:", foundUser.ID)
			writeErrorResponse(http.StatusInternalServerError, resWriter, "Something went wrong")
			return
		}

		// Create and Marshal response body
		responseBody := createUserResponse{
			ID:        foundUser.ID,
			CreatedAt: foundUser.CreatedAt,
			UpdatedAt: foundUser.UpdatedAt,
			Email:     foundUser.Email,
			Token:     jwtToken,
		}

		responseBodyMarshalled, err := json.Marshal(responseBody)
		if err != nil {
			fmt.Println("Failed to marshal response: ", err)
			writeErrorResponse(http.StatusInternalServerError, resWriter, "Something went wrong")
			return
		}

		// Set header and status code
		resWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
		resWriter.WriteHeader(http.StatusOK)

		// Set response body
		resWriter.Write(responseBodyMarshalled)
	}
}
