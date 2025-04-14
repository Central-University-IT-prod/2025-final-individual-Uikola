package mlscore

import "api/internal/repository"

// Usecase represents the business logic layer for ml-score operations.
type Usecase struct {
	mlScoreRepository repository.MLScoreRepository
}

// NewUsecase initializes a new Usecase instance.
func NewUsecase(mlScoreRepository repository.MLScoreRepository) *Usecase {
	return &Usecase{
		mlScoreRepository: mlScoreRepository,
	}
}
