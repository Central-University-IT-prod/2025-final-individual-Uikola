package campaign

import (
	"api/pkg/s3"
	"context"
)

// UploadImage uploads an image to S3, updates the campaign with the image ID and returns the image URL.
func (uc *Usecase) UploadImage(ctx context.Context, fileData s3.FileDataType, advertiserID, campaignID string) (string, error) {
	_, err := uc.advertiserRepository.GetByID(ctx, advertiserID)
	if err != nil {
		return "", err
	}

	campaign, err := uc.campaignRepository.Get(ctx, advertiserID, campaignID)
	if err != nil {
		return "", err
	}

	link, objectID, err := uc.s3Client.CreateOne(ctx, fileData, "image/png")
	if err != nil {
		return "", err
	}

	campaign.ImageID = objectID
	err = uc.campaignRepository.Update(ctx, campaign)
	if err != nil {
		return "", err
	}

	return link, nil
}
