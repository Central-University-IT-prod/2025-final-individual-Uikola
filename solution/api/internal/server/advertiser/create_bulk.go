package advertiser

import (
	"api/internal/entity/request"
	"api/internal/errorz"
	"api/internal/server/helper"
	"encoding/json"
	"net/http"
)

// CreateBulk handles bulk advertiser creation requests.
func (h *Handler) CreateBulk(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var requestBody []request.Advertiser
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return errorz.APIError{
			Status: http.StatusBadRequest,
			Err:    err,
			Msg:    "failed to decode request body",
		}
	}

	for _, advertiser := range requestBody {
		if err := advertiser.Validate(); err != nil {
			return errorz.APIError{
				Status: http.StatusBadRequest,
				Err:    err,
				Msg:    "failed to validate advertiser",
			}
		}
	}

	advertisers, err := h.advertiserUsecase.CreateBulk(ctx, requestBody)
	switch {
	case err != nil:
		return errorz.APIError{
			Status: http.StatusInternalServerError,
			Err:    err,
			Msg:    "failed to create advertisers bulk",
		}
	}

	return helper.WriteJSON(w, http.StatusOK, advertisers)
}
