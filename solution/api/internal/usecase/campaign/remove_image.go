package campaign

import "context"

// RemoveImage removes the image from the campaign by clearing the ImageID field.
func (uc *Usecase) RemoveImage(ctx context.Context, advertiserID, campaignID string) error {
	_, err := uc.advertiserRepository.GetByID(ctx, advertiserID)
	if err != nil {
		return err
	}

	campaign, err := uc.campaignRepository.Get(ctx, advertiserID, campaignID)
	if err != nil {
		return err
	}

	err = uc.s3Client.DeleteOne(ctx, campaign.CampaignID)
	if err != nil {
		return err
	}

	campaign.ImageID = ""
	return uc.campaignRepository.Update(ctx, campaign)
}
