package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

func logRequestHeaders(r *http.Request) {
	// Marshal the r.Header map into indented JSON
	prettyHeaders, err := json.MarshalIndent(r.Header, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling headers:", err)
		return
	}

	// Print JSON string
	fmt.Printf("--- Pretty JSON Headers ---\n%s\n", string(prettyHeaders))
}

func logJSONRequestBody(r *http.Request) {
	// 1. Read the body bytes safely
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Error reading body: %v\n", err)
		return
	}

	// 2. Restore r.Body so it can be read again downstream
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// 3. Pretty print the JSON payload
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, bodyBytes, "", "  "); err == nil {
		fmt.Printf("--- Request Body ---\n%s\n--------------------\n", prettyJSON.String())
	} else {
		// Fallback to raw string if it's not valid JSON
		fmt.Printf("--- Raw Request Body ---\n%s\n--------------------\n", string(bodyBytes))
	}
}
