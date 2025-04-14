package campaign

import (
	"api/internal/errorz"
	"api/internal/server/helper"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// RemoveImage handles the HTTP request to remove an image from a specific campaign.
func (h *Handler) RemoveImage(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	advertiserID := chi.URLParam(r, "advertiserId")
	valid := helper.IsValidUUID(advertiserID)
	if !valid {
		return errorz.APIError{
			Status: http.StatusBadRequest,
			Err:    errorz.ErrInvalidAdvertiserID,
			Msg:    "invalid advertiser id",
		}
	}

	campaignID := chi.URLParam(r, "campaignId")
	valid = helper.IsValidUUID(campaignID)
	if !valid {
		return errorz.APIError{
			Status: http.StatusBadRequest,
			Err:    errorz.ErrInvalidCampaignID,
			Msg:    "invalid campaign id",
		}
	}

	err := h.campaignUsecase.RemoveImage(ctx, advertiserID, campaignID)
	switch {
	case errors.Is(err, errorz.ErrAdvertiserNotFound):
		return errorz.APIError{
			Status: http.StatusNotFound,
			Err:    err,
			Msg:    "advertiser not found",
		}
	case errors.Is(err, errorz.ErrCampaignNotFound):
		return errorz.APIError{
			Status: http.StatusNotFound,
			Err:    err,
			Msg:    "campaign not found",
		}
	case err != nil:
		return errorz.APIError{
			Status: http.StatusBadRequest,
			Err:    err,
			Msg:    "failed to remove image",
		}
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
