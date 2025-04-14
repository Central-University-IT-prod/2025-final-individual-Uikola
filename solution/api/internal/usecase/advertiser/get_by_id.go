package advertiser

import (
	"api/internal/entity/response"
	"context"
)

// GetByID retrieves an advertiser by their advertiser ID.
func (uc *Usecase) GetByID(ctx context.Context, advertiserID string) (response.Advertiser, error) {
	advertiser, err := uc.advertiserRepository.GetByID(ctx, advertiserID)
	if err != nil {
		return response.Advertiser{}, err
	}

	return advertiser.ToResponse(), err
}
