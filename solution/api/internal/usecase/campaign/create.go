package campaign

import (
	"api/internal/entity/request"
	"api/internal/entity/response"
	"context"
)

// Create processes a campaign creation request.
func (uc *Usecase) Create(ctx context.Context, req request.CreateCampaign, advertiserID string) (response.Campaign, error) {
	_, err := uc.advertiserRepository.GetByID(ctx, advertiserID)
	if err != nil {
		return response.Campaign{}, err
	}

	campaign, err := uc.campaignRepository.Create(ctx, req.ToCampaign(advertiserID))
	if err != nil {
		return response.Campaign{}, err
	}

	return campaign.ToResponse(""), nil
}
