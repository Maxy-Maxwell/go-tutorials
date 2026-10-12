package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func writeResponseBody(statusCode int, resWriter http.ResponseWriter, body interface{}) {
	responseBody, err := json.Marshal(body)

	if err != nil {
		fmt.Println("Failed to marshal error response:", err)
		writeErrorResponse(500, resWriter, "Something went wrong")
		return
	}

	resWriter.WriteHeader(statusCode)
	resWriter.Write(responseBody)
}

func writeErrorResponse(statusCode int, resWriter http.ResponseWriter, message string) {
	// We set the status code first in case writing to the response body also fails
	resWriter.WriteHeader(statusCode)

	type errorBody struct {
		Error string `json:"error"`
	}

	responseBody, err := json.Marshal(errorBody{
		Error: fmt.Sprintf("%v", message),
	})

	if err != nil {
		fmt.Println("Failed to marshal error response:", err)
		return
	}

	resWriter.Write(responseBody)
}
