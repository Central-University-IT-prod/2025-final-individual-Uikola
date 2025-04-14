package client

import "api/internal/usecase"

// Handler is the HTTP handler for client-related operations.
type Handler struct {
	clientUsecase usecase.ClientUsecase
}

// NewHandler initializes a new Handler instance.
func NewHandler(clientUsecase usecase.ClientUsecase) *Handler {
	return &Handler{
		clientUsecase: clientUsecase,
	}
}
