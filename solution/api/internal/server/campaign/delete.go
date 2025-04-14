package campaign

import (
	"api/internal/errorz"
	"api/internal/server/helper"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Delete handles campaign deletion request.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) error {
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

	err := h.campaignUsecase.Delete(ctx, advertiserID, campaignID)
	switch {
	case errors.Is(err, errorz.ErrCampaignNotFound):
		w.WriteHeader(http.StatusNoContent)
		return nil
	case err != nil:
		return errorz.APIError{
			Status: http.StatusInternalServerError,
			Err:    err,
			Msg:    "failed to delete advertiser campaign",
		}
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
