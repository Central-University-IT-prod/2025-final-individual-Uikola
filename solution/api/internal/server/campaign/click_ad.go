package campaign

import (
	"api/internal/entity/request"
	"api/internal/errorz"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) ClickAd(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	adID := chi.URLParam(r, "adId")

	var requestBody request.ClickAd
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

	err := h.campaignUsecase.ClickAd(ctx, requestBody, adID)
	switch {
	case errors.Is(err, errorz.ErrCampaignNotFound):
		return errorz.APIError{
			Status: http.StatusNotFound,
			Err:    err,
			Msg:    "campaign not found",
		}
	case errors.Is(err, errorz.ErrImpressionNotFound):
		return errorz.APIError{
			Status: http.StatusForbidden,
			Err:    err,
			Msg:    "impression not found",
		}
	case err != nil:
		return errorz.APIError{
			Status: http.StatusInternalServerError,
			Err:    err,
			Msg:    "failed to click ad",
		}
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
