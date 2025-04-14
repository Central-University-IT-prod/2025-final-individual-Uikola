package campaign

import (
	"api/internal/errorz"
	"api/internal/server/helper"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// ListWithPagination handles an HTTP request to fetch a list of campaigns by advertiser ID.
func (h *Handler) ListWithPagination(w http.ResponseWriter, r *http.Request) error {
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

	var size, page int
	size, err := strconv.Atoi(r.URL.Query().Get("size"))
	if err != nil {
		size = 10
	}

	page, err = strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	campaigns, count, err := h.campaignUsecase.ListWithPagination(ctx, advertiserID, size, page)
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
			Msg:    "failed to retrieve a list of campaigns",
		}
	}

	w.Header().Set("X-Total-Count", strconv.Itoa(int(count)))
	return helper.WriteJSON(w, http.StatusOK, campaigns)
}
