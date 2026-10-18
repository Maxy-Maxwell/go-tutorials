package main

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/google/uuid"

	"github.com/Maxy-Maxwell/go-tutorials/chirpy/internal/auth"
	"github.com/Maxy-Maxwell/go-tutorials/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	jwt_secret     string
	polka_apikey   string
}

func initializeApiConfig(dbQueries *database.Queries, platform string, jwt_secret string, polka_apikey string) apiConfig {
	return apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
		jwt_secret:     jwt_secret,
		polka_apikey:   polka_apikey,
	}
}

func (conf *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	// Anything that happens here, is run when the hanlder is first initialized (in setupRoutes)

	return http.HandlerFunc(func(resWriter http.ResponseWriter, request *http.Request) {
		// Anything here happens whenever the handler is called (aka calls to the /app/)
		conf.fileserverHits.Add(1)
		next.ServeHTTP(resWriter, request)
	})
}

func (conf *apiConfig) middlewareJWTAuth(
	handler func(conf *apiConfig, userId uuid.UUID, resWriter http.ResponseWriter, request *http.Request)) func(resWriter http.ResponseWriter, request *http.Request) {

	return func(resWriter http.ResponseWriter, request *http.Request) {
		fmt.Println("Inside JWT middleware...")
		token, err := auth.GetBearerToken(request.Header)
		if err != nil {
			writeResponseBody(http.StatusUnauthorized, resWriter, "Unauthorized")
			return
		}

		userId, err := auth.ValidateJWT(token, conf.jwt_secret)
		if err != nil {
			writeResponseBody(http.StatusUnauthorized, resWriter, "Unauthorized")
			return
		}

		fmt.Println("Calling Handler...")

		handler(conf, userId, resWriter, request)
	}
}
