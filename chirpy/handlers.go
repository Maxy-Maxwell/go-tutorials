package main

import (
	"fmt"
	"io"
	"net/http"
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
