package middleware

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	StatusCode int         `json:"statusCode"`
	Error      string      `json:"error"`
	Details    interface{} `json:"details,omitempty"`
}

func ErrorHandler(w http.ResponseWriter, statusCode int, err error, details interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := ErrorResponse{
		StatusCode: statusCode,
		Error:      err.Error(),
		Details:    details,
	}

	json.NewEncoder(w).Encode(resp)
}
