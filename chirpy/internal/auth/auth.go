package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		fmt.Println("Error when hashing:", err)
		return "", err
	}

	return hash, nil
}

func CheckPasswordHash(password string, hash string) (bool, error) {
	doesMatch, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		fmt.Println("Error while comparing hash:", err)
		return false, err
	}

	return doesMatch, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		fmt.Println("Not Auth header found")
		return "", fmt.Errorf("No Auth header found")
	}

	token, hasPrefix := strings.CutPrefix(authHeader, "Bearer ")
	if !hasPrefix {
		fmt.Println("Malformed bearer token: ", authHeader)
		return "", fmt.Errorf("Malformed Auth Header.")
	}

	return token, nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})

	fmt.Println("expiresIn:", expiresIn)
	fmt.Println("Token expiry: ", jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)))
	fmt.Println("Now: ", time.Now().UTC())

	signedToken, err := token.SignedString([]byte(tokenSecret))

	if err != nil {
		fmt.Println("Error signing string: ", err)
		return "", err
	}

	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	resultToken, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(tokenSecret), nil
		},
	)
	if err != nil {
		fmt.Println("Token: ", resultToken)
		fmt.Println("Error Parsing Claims:", err)
		return uuid.UUID{}, err
	}

	claims, ok := resultToken.Claims.(*jwt.RegisteredClaims)
	if !ok {
		fmt.Println("Incorrect claims type, expected RegisteredClaims", resultToken.Claims)
		return uuid.UUID{}, fmt.Errorf("Invalid token claims - incorrect type")
	}

	if claims.ExpiresAt.Before(time.Now().UTC()) {
		fmt.Println("Expired Token:", claims.ExpiresAt)
		return uuid.UUID{}, fmt.Errorf("Expired token")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		fmt.Println("Failed to parse userID from token:", claims.Subject)
		return uuid.UUID{}, err
	}

	return userID, nil
}

func MakeRefreshToken() string {
	token := make([]byte, 32)
	rand.Read(token)
	return hex.EncodeToString(token)
}
