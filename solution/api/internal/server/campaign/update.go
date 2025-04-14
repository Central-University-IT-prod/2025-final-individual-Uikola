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

// Update handles an HTTP request to update campaign data.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) error {
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

	var requestBody request.UpdateCampaign
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

	campaign, err := h.campaignUsecase.Update(ctx, advertiserID, campaignID, requestBody)
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
	case errors.Is(err, errorz.ErrCampaignIsOver):
		return errorz.APIError{
			Status: http.StatusBadRequest,
			Err:    err,
			Msg:    "campaign is over",
		}
	case errors.Is(err, errorz.ErrCampaignIsActive):
		return errorz.APIError{
			Status: http.StatusBadRequest,
			Err:    err,
			Msg:    "you can not update these fields, campaign is active",
		}
	case errors.Is(err, errorz.ErrInvalidDate):
		return errorz.APIError{
			Status: http.StatusBadRequest,
			Err:    err,
			Msg:    "invalid date",
		}
	case err != nil:
		return errorz.APIError{
			Status: http.StatusInternalServerError,
			Err:    err,
			Msg:    "failed to update campaign",
		}
	}

	return helper.WriteJSON(w, http.StatusOK, campaign)
}
