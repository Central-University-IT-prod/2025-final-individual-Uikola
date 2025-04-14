package helper

import (
	"encoding/json"
	"net/http"
)

// WriteJSON writes a JSON response with the given status code and data.
func WriteJSON(w http.ResponseWriter, code int, data any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(data)
}
