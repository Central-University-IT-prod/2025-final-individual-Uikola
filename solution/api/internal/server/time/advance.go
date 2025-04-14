package time

import (
	"api/internal/entity/request"
	"api/internal/entity/response"
	"api/internal/errorz"
	"api/internal/server/helper"
	"encoding/json"
	"net/http"
)

// SetDate handles the HTTP request to update the current system date.
func (h *Handler) SetDate(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var requestBody request.Advance
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return errorz.APIError{
			Status: http.StatusBadRequest,
			Err:    err,
			Msg:    "failed to decode request body",
		}
	}

	if err := requestBody.Validate(); err != nil {
		return errorz.APIError{
			Status: http.StatusBadRequest,
			Err:    err,
			Msg:    "failed to validate request body",
		}
	}

	h.timeRepository.Set(ctx, requestBody.CurrentDate, 0)

	return helper.WriteJSON(w, http.StatusOK, response.Advance{CurrentDate: requestBody.CurrentDate})
}
