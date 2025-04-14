package campaign

import "context"

// Delete removes a campaign by its AdvertiserID and CampaignID.
func (uc *Usecase) Delete(ctx context.Context, advertiserID, campaignID string) error {
	campaign, err := uc.campaignRepository.Get(ctx, advertiserID, campaignID)
	if err != nil {
		return err
	}

	err = uc.s3Client.DeleteOne(ctx, campaign.ImageID)
	if err != nil {
		return err
	}

	err = uc.statisticRepository.Delete(ctx, advertiserID, campaignID)
	if err != nil {
		return err
	}

	return uc.campaignRepository.Delete(ctx, advertiserID, campaignID)
}
