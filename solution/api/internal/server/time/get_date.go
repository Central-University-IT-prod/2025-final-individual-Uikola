package time

import (
	"api/internal/entity/response"
	"api/internal/errorz"
	"api/internal/server/helper"
	"net/http"
)

func (h *Handler) GetDate(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	date, err := h.timeRepository.Get(ctx)
	if err != nil {
		return errorz.APIError{
			Status: http.StatusInternalServerError,
			Err:    err,
			Msg:    "failed to get date",
		}
	}

	return helper.WriteJSON(w, http.StatusOK, response.Advance{CurrentDate: date})
}
