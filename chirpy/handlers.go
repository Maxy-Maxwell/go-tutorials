package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	currentTime := time.Now().UTC()

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
			fmt.Println("No Chirp id found:", chirpId)
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

func deleteChirpHandler(conf *apiConfig, headerUserId uuid.UUID, resWriter http.ResponseWriter, request *http.Request) {
	chirpId := request.PathValue("id")

	if chirpId == "" {
		fmt.Println("No Chirp id found:", chirpId)
		writeErrorResponse(http.StatusInternalServerError, resWriter, "Something went wrong")
		return
	}

	parsedChirpId, err := uuid.Parse(chirpId)
	if err != nil {
		fmt.Println("Failed to parse chirpId into UUID: ", chirpId)
		writeErrorResponse(http.StatusInternalServerError, resWriter, "Something went wrong")
		return
	}

	fmt.Printf("Deleting chirp with id: %v...\n", parsedChirpId)

	// We only do this to present nice error messages to the user. We could just do the delete and it'll return nil error even if the chirp doesn't exist or is owned by another user
	queriedChirp, err := conf.db.GetChirpById(request.Context(), parsedChirpId)
	if err != nil {
		fmt.Println("Failed to get chirp:", err)
		writeErrorResponse(http.StatusNotFound, resWriter, "Not Found")
		return
	}

	if queriedChirp.UserID != headerUserId {
		fmt.Printf("Attempted to delete another user's chirp. %v (chirp owner), %v (current user)\n", queriedChirp.UserID, headerUserId)
		writeErrorResponse(http.StatusForbidden, resWriter, "Forbidden")
		return
	}

	// Delete the chirp
	err = conf.db.DeleteChirpById(request.Context(), database.DeleteChirpByIdParams{
		ID:     parsedChirpId,
		UserID: headerUserId,
	})
	if err != nil {
		fmt.Println("Failed to delete chirp:", err)
		writeErrorResponse(http.StatusInternalServerError, resWriter, "Something went wrong")
		return
	}

	// Set headers
	resWriter.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Set status code
	resWriter.WriteHeader(http.StatusNoContent)
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
		currentTime := time.Now().UTC()
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
		type loginRequest struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		type loginResponse struct {
			ID           uuid.UUID `json:"id"`
			CreatedAt    time.Time `json:"created_at"`
			UpdatedAt    time.Time `json:"updated_at"`
			Email        string    `json:"email"`
			Token        string    `json:"token"`
			RefreshToken string    `json:"refresh_token"`
		}

		// Get Request Body
		decoder := json.NewDecoder(request.Body)
		parsedReq := loginRequest{}

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
		jwtToken, err := auth.MakeJWT(foundUser.ID, conf.jwt_secret, time.Duration(int(time.Hour.Seconds()))*time.Second)
		if err != nil {
			fmt.Println("Failed to generate token: ", err, ". UserId:", foundUser.ID)
			writeErrorResponse(http.StatusInternalServerError, resWriter, "Something went wrong")
			return
		}

		// Get and store refresh token
		refreshToken := auth.MakeRefreshToken()
		currentTime := time.Now().UTC()

		_, err = conf.db.CreateRefreshToken(request.Context(), database.CreateRefreshTokenParams{
			Token:     refreshToken,
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
			UserID:    foundUser.ID,
			ExpiresAt: time.Now().UTC().Add(60 * 24 * time.Hour),
		})

		if err != nil {
			fmt.Println("Failed to create refresh token:", err)
			writeErrorResponse(http.StatusInternalServerError, resWriter, "Something went wrong")
			return
		}

		// Create and Marshal response body
		responseBody := loginResponse{
			ID:           foundUser.ID,
			CreatedAt:    foundUser.CreatedAt,
			UpdatedAt:    foundUser.UpdatedAt,
			Email:        foundUser.Email,
			Token:        jwtToken,
			RefreshToken: refreshToken,
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

func (conf *apiConfig) createRefreshTokenHandler() func(resWriter http.ResponseWriter, request *http.Request) {
	return func(resWriter http.ResponseWriter, request *http.Request) {
		type createRefreshTokenResponse struct {
			Token string `json:"token"`
		}

		// Get refresh token from header
		refreshToken, err := auth.GetBearerToken(request.Header)
		if err != nil {
			fmt.Println("Failed to get auth header's bearer token:", err)
			writeErrorResponse(http.StatusUnauthorized, resWriter, "Unauthorized")
			return
		}

		// Retrieve token from database
		retrievedToken, err := conf.db.GetRefreshTokenByToken(request.Context(), refreshToken)
		if err != nil {
			fmt.Println("Failed to query refresh token:", err)
			writeErrorResponse(http.StatusUnauthorized, resWriter, "Unauthorized")
			return
		}

		currentTime := time.Now().UTC()

		// Ensure the token isn't revoked
		if retrievedToken.RevokedAt.Valid && retrievedToken.RevokedAt.Time.Before(currentTime) {
			fmt.Println("Revoked token: ", retrievedToken.RevokedAt.Time)
			writeErrorResponse(http.StatusUnauthorized, resWriter, "Unauthorized")
			return
		}

		// Ensure token hasn't expired
		if retrievedToken.ExpiresAt.Before(currentTime) {
			fmt.Printf("Expired token. Expired at: %v. Time now: %v.\n", retrievedToken.ExpiresAt, time.Now().UTC())
			writeErrorResponse(http.StatusUnauthorized, resWriter, "Unauthorized")
			return
		}

		// Grab a new token for em
		jwtToken, err := auth.MakeJWT(retrievedToken.UserID, conf.jwt_secret, time.Duration(int(time.Hour.Seconds()))*time.Second)
		if err != nil {
			fmt.Println("Failed to generate token: ", err, ". UserId:", retrievedToken.UserID)
			writeErrorResponse(http.StatusInternalServerError, resWriter, "Something went wrong")
			return
		}

		// Create and Marshal response body
		responseBody := createRefreshTokenResponse{
			Token: jwtToken,
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

func (conf *apiConfig) createRevokeTokenHandler() func(resWriter http.ResponseWriter, request *http.Request) {
	return func(resWriter http.ResponseWriter, request *http.Request) {
		// Get refresh token from header
		refreshToken, err := auth.GetBearerToken(request.Header)
		if err != nil {
			fmt.Println("Failed to get auth header's bearer token:", err)
			writeErrorResponse(http.StatusUnauthorized, resWriter, "Unauthorized")
			return
		}

		revokedAt := sql.NullTime{
			Time:  time.Now().UTC().Add(-1 * time.Second), // "-1" Not really needed, but makes me feel better when automated tests run
			Valid: true,
		}

		// Update the refresh_tokens table to revoke token
		_, err = conf.db.UpdateRefreshTokenByToken(
			request.Context(),
			database.UpdateRefreshTokenByTokenParams{
				Token:     refreshToken,
				UpdatedAt: time.Now().UTC(),
				RevokedAt: revokedAt,
			},
		)

		if err != nil {
			fmt.Println("Error when updating refresh token: ", err)
			writeErrorResponse(http.StatusUnauthorized, resWriter, "Unauthorized")
			return
		}

		// Set header and status code
		resWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
		resWriter.WriteHeader(http.StatusNoContent)
	}
}

func createUpdateUserHandler(conf *apiConfig, headerUserId uuid.UUID, resWriter http.ResponseWriter, request *http.Request) {
	type updateUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type updateUserResponse struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	decoder := json.NewDecoder(request.Body)
	parsedReq := updateUserRequest{}

	if err := decoder.Decode(&parsedReq); err != nil {
		fmt.Println("Failed to parse request: ", err)
		writeErrorResponse(http.StatusInternalServerError, resWriter, "Something went wrong")
		return
	}

	// Make sure the password and email are not null
	if parsedReq.Email == "" || parsedReq.Password == "" {
		fmt.Println("Invalid request body. Email and Password are required.")
		writeErrorResponse(http.StatusBadRequest, resWriter, "Invalid request body")
		return
	}

	hashedPassword, err := auth.HashPassword(parsedReq.Password)
	if err != nil {
		fmt.Println("Failed to hash password: ", err)
		writeErrorResponse(http.StatusInternalServerError, resWriter, "Something went wrong")
		return
	}

	// Get the user from the header UserID. Don't actually need this, as we can probably enforce the Id
	updatedUser, err := conf.db.UpdateUserById(request.Context(), database.UpdateUserByIdParams{
		ID:             headerUserId,
		Email:          parsedReq.Email,
		HashedPassword: hashedPassword,
		UpdatedAt:      time.Now().UTC(),
	})
	if err != nil {
		fmt.Println("Failed to update user by id: ", err)
		writeErrorResponse(http.StatusInternalServerError, resWriter, "Something went wrong")
		return
	}

	writeResponseBody(http.StatusOK, resWriter, updateUserResponse{
		ID:        headerUserId,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
		Email:     updatedUser.Email,
	})
}
