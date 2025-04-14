package campaign

import (
	"context"
)

// Moderate reviews the advertising campaign.
func (uc *Usecase) Moderate(ctx context.Context, campaignID string, passedModeration bool) error {
	campaign, err := uc.campaignRepository.GetByID(ctx, campaignID)
	if err != nil {
		return err
	}

	if passedModeration {
		campaign.PassedModeration = true
		return uc.campaignRepository.Update(ctx, campaign)
	}

	return uc.Delete(ctx, campaign.AdvertiserID, campaignID)
}
