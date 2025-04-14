package mlscore

import (
	"api/internal/entity/request"
	"context"
)

// Create processes a ml-score creation request.
func (uc *Usecase) Create(ctx context.Context, req request.MLScore) error {
	return uc.mlScoreRepository.Create(ctx, req.ToMLScore())
}
