package domain

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type NotFoundError struct {
	Message string `json:"error"`
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Not Found: %s", e.Message)
}

func ErrorResponseF(w http.ResponseWriter, module string, statusCode int, message string) {
	log.Printf("%s %s", module, message)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ErrorResponse{
		Message: message,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}
