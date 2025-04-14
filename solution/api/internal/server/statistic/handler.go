package statistic

import "api/internal/usecase"

// Handler is the HTTP handler for statistic-related operations.
type Handler struct {
	statisticUsecase usecase.StatisticUsecase
}

// NewHandler initializes a new Handler instance.
func NewHandler(statisticUsecase usecase.StatisticUsecase) *Handler {
	return &Handler{
		statisticUsecase: statisticUsecase,
	}
}
