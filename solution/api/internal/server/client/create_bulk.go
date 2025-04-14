package client

import (
	"api/internal/entity/request"
	"api/internal/errorz"
	"api/internal/server/helper"
	"encoding/json"
	"net/http"
)

// CreateBulk handles bulk client creation requests.
func (h *Handler) CreateBulk(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var requestBody []request.Client
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return errorz.APIError{
			Status: http.StatusBadRequest,
			Err:    err,
			Msg:    "failed to decode request body",
		}
	}

	for _, client := range requestBody {
		if err := client.Validate(); err != nil {
			return errorz.APIError{
				Status: http.StatusBadRequest,
				Err:    err,
				Msg:    "failed to validate client",
			}
		}
	}

	clients, err := h.clientUsecase.CreateBulk(ctx, requestBody)
	switch {
	case err != nil:
		return errorz.APIError{
			Status: http.StatusInternalServerError,
			Err:    err,
			Msg:    "failed to create clients bulk",
		}
	}

	return helper.WriteJSON(w, http.StatusOK, clients)
}
