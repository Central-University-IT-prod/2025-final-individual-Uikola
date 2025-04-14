package helper

import (
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"

	"api/internal/entity/response"

	"api/internal/errorz"
)

// MakeHandler converts a function that returns an error into a http.HandlerFunc.
// It wraps the function with error handling logic, ensuring proper logging and HTTP response formatting.
func MakeHandler(fn func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			var e errorz.APIError
			if errors.As(err, &e) {
				log.Error().Msg(err.Error())
				_ = WriteJSON(w, e.Status, response.Error{Error: e.Msg})
			}
		}
	}
}
