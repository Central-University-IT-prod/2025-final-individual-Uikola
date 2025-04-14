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

// GenerateAdText handles the generation of ad text using an AI model.
func (h *Handler) GenerateAdText(w http.ResponseWriter, r *http.Request) error {
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

	var requestBody request.GenerateAdText
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

	resp, err := h.campaignUsecase.GenerateText(ctx, requestBody, advertiserID)
	switch {
	case errors.Is(err, errorz.ErrAdvertiserNotFound):
		return errorz.APIError{
			Status: http.StatusNotFound,
			Err:    err,
			Msg:    "advertiser not found",
		}
	case err != nil:
		return errorz.APIError{
			Status: http.StatusInternalServerError,
			Err:    err,
			Msg:    "failed to generate ad text",
		}
	}

	return helper.WriteJSON(w, http.StatusOK, resp)
}
