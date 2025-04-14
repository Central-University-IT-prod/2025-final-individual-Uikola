package advertiser

import "api/internal/usecase"

// Handler is the HTTP handler for advertiser-related operations.
type Handler struct {
	advertiserUsecase usecase.AdvertiserUsecase
}

// NewHandler initializes a new Handler instance.
func NewHandler(advertiserUsecase usecase.AdvertiserUsecase) *Handler {
	return &Handler{
		advertiserUsecase: advertiserUsecase,
	}
}
