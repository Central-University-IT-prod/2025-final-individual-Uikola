package campaign

import (
	"api/internal/repository"
	"api/internal/usecase"
)

// Handler is the HTTP handler for campaign-related operations.
type Handler struct {
	campaignUsecase usecase.CampaignUsecase

	timeRepository repository.TimeRepository
}

// NewHandler initializes a new Handler instance.
func NewHandler(campaignUsecase usecase.CampaignUsecase, timeRepository repository.TimeRepository) *Handler {
	return &Handler{
		campaignUsecase: campaignUsecase,

		timeRepository: timeRepository,
	}
}
