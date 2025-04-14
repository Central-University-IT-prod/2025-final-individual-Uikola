package campaign

import (
	"api/internal/entity/response"
	"context"
)

// GetByID retrieves a campaign by its advertiser and campaign id.
func (uc *Usecase) GetByID(ctx context.Context, advertiserID, campaignID string) (response.Campaign, error) {
	_, err := uc.advertiserRepository.GetByID(ctx, advertiserID)
	if err != nil {
		return response.Campaign{}, err
	}

	campaign, err := uc.campaignRepository.Get(ctx, advertiserID, campaignID)
	if err != nil {
		return response.Campaign{}, err
	}

	link, err := uc.s3Client.GetOne(ctx, campaign.ImageID)
	if err != nil {
		return response.Campaign{}, err
	}

	return campaign.ToResponse(link), nil
}
