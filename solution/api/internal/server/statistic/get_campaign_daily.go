package statistic

import (
	"api/internal/errorz"
	"api/internal/server/helper"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// GetCampaignDaily handles the HTTP request to retrieve daily statistics for a specific campaign.
func (h *Handler) GetCampaignDaily(w http.ResponseWriter, r *http.Request) error {
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

	statistics, err := h.statisticUsecase.GetCampaignDaily(ctx, campaignID)
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
			Msg:    "failed to get campaign daily statistic",
		}
	}

	return helper.WriteJSON(w, http.StatusOK, statistics)
}
