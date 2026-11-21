// Package response provides small helpers for HTTP responses.
//
// Responsibilities:
//   - Standardize JSON responses across handler
//   - Reduce duplication of response-writing logic
//   - Ensure consistent headers and error formats
//
// This package does NOT:
//   - Contain business logic
//   - Know anything about authentication or routing
//
// All handler are encouraged to use these helpers
// to keep responses consistent across the API.

// JSON sends a JSON response with status code
package response

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Error sends a JSON error response
func Error(w http.ResponseWriter, status int, msg string) {
	JSON(w, status, map[string]string{"error": msg})
}
