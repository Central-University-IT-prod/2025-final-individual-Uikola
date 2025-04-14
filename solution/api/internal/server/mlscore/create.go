package mlscore

import (
	"api/internal/entity/request"
	"api/internal/errorz"
	"encoding/json"
	"net/http"
)

// Create handles ml-score creation request.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var requestBody request.MLScore
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

	err := h.mlScoreUsecase.Create(ctx, requestBody)
	switch {
	case err != nil:
		return errorz.APIError{
			Status: http.StatusInternalServerError,
			Err:    err,
			Msg:    "failed to create ml-score",
		}
	}

	w.WriteHeader(http.StatusOK)
	return nil
}
