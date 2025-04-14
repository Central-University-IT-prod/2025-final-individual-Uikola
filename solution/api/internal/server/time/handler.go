package time

import (
	"api/internal/repository"
)

// Handler is the HTTP handler for time-related operations.
type Handler struct {
	timeRepository repository.TimeRepository
}

// NewHandler initializes a new Handler instance.
func NewHandler(timeRepository repository.TimeRepository) *Handler {
	return &Handler{
		timeRepository: timeRepository,
	}
}
