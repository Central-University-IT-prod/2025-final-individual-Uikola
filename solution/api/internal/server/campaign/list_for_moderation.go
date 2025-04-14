package campaign

import (
	"api/internal/errorz"
	"api/internal/server/helper"
	"net/http"
	"strconv"
)

// ListForModeration handles an HTTP request to fetch a list of campaigns for moderation.
func (h *Handler) ListForModeration(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var size, page int
	size, err := strconv.Atoi(r.URL.Query().Get("size"))
	if err != nil {
		size = 10
	}

	page, err = strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	campaigns, count, err := h.campaignUsecase.ListForModerationWithPagination(ctx, size, page)
	switch {
	case err != nil:
		return errorz.APIError{
			Status: http.StatusInternalServerError,
			Err:    err,
			Msg:    "failed to retrieve a list of campaigns for moderation",
		}
	}

	w.Header().Set("X-Total-Count", strconv.Itoa(int(count)))
	return helper.WriteJSON(w, http.StatusOK, campaigns)
}
