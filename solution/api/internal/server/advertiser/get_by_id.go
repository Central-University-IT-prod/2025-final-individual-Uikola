package advertiser

import (
	"api/internal/errorz"
	"api/internal/server/helper"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// GetByID handles an HTTP request to fetch an advertiser by its ID.
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) error {
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

	advertiser, err := h.advertiserUsecase.GetByID(ctx, advertiserID)
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
			Msg:    "failed to get advertiser by id",
		}
	}

	return helper.WriteJSON(w, http.StatusOK, advertiser)
}
