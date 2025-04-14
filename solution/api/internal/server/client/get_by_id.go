package client

import (
	"api/internal/errorz"
	"api/internal/server/helper"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// GetByID handles an HTTP request to fetch a client by its ID.
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	clientID := chi.URLParam(r, "clientId")
	valid := helper.IsValidUUID(clientID)
	if !valid {
		return errorz.APIError{
			Status: http.StatusBadRequest,
			Err:    errorz.ErrInvalidClientID,
			Msg:    "invalid client id",
		}
	}

	client, err := h.clientUsecase.GetByID(ctx, clientID)
	switch {
	case errors.Is(err, errorz.ErrClientNotFound):
		return errorz.APIError{
			Status: http.StatusNotFound,
			Err:    errorz.ErrClientNotFound,
			Msg:    "client not found",
		}
	case err != nil:
		return errorz.APIError{
			Status: http.StatusInternalServerError,
			Err:    err,
			Msg:    "failed to get client by id",
		}
	}

	return helper.WriteJSON(w, http.StatusOK, client)
}
