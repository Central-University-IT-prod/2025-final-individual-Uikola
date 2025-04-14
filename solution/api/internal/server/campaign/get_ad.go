package campaign

import (
	"api/internal/errorz"
	"api/internal/server/helper"
	"errors"
	"net/http"
)

// GetAdd handles the HTTP request to retrieve the most suitable advertisement for a client.
func (h *Handler) GetAdd(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	clientID := r.URL.Query().Get("client_id")
	valid := helper.IsValidUUID(clientID)
	if !valid {
		return errorz.APIError{
			Status: http.StatusBadRequest,
			Err:    errorz.ErrInvalidClientID,
			Msg:    "invalid client id",
		}
	}

	add, err := h.campaignUsecase.GetAd(ctx, clientID)
	switch {
	case errors.Is(err, errorz.ErrClientNotFound):
		return errorz.APIError{
			Status: http.StatusNotFound,
			Err:    err,
			Msg:    "client not found",
		}
	case errors.Is(err, errorz.ErrNoCampaignsFound):
		return errorz.APIError{
			Status: http.StatusNotFound,
			Err:    err,
			Msg:    "no campaigns found",
		}
	case err != nil:
		return errorz.APIError{
			Status: http.StatusInternalServerError,
			Err:    err,
			Msg:    "failed to get add",
		}
	}

	return helper.WriteJSON(w, http.StatusOK, add)
}
