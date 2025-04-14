package campaign

import (
	"api/internal/entity/request"
	"api/internal/errorz"
	"api/internal/server/helper"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Moderate handles the moderation of an advertising campaign.
func (h *Handler) Moderate(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	campaignID := chi.URLParam(r, "campaignId")
	valid := helper.IsValidUUID(campaignID)
	if !valid {
		return errorz.APIError{
			Status: http.StatusBadRequest,
			Err:    errorz.ErrInvalidCampaignID,
			Msg:    "invalid campaign id",
		}
	}

	var requestBody request.Moderate
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

	err := h.campaignUsecase.Moderate(ctx, campaignID, *requestBody.PassedModeration)
	switch {
	case errors.Is(err, errorz.ErrCampaignNotFound):
		return errorz.APIError{
			Status: http.StatusNotFound,
			Err:    err,
			Msg:    "campaign not found",
		}
	case err != nil:
		return errorz.APIError{
			Status: http.StatusInternalServerError,
			Err:    err,
			Msg:    "failed tp moderate campaign",
		}
	}

	w.WriteHeader(http.StatusOK)
	return nil
}
