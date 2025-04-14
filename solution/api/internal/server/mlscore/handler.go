package mlscore

import "api/internal/usecase"

// Handler is the HTTP handler for ml-score-related operations.
type Handler struct {
	mlScoreUsecase usecase.MLScoreUsecase
}

// NewHandler initializes a new Handler instance.
func NewHandler(mlScoreUsecase usecase.MLScoreUsecase) *Handler {
	return &Handler{
		mlScoreUsecase: mlScoreUsecase,
	}
}
