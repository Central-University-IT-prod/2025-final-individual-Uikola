package campaign

import (
	"api/internal/entity/response"
	"context"
)

// ListWithPagination retrieves a list of campaigns by its advertiser and campaign id with pagination.
func (uc *Usecase) ListWithPagination(ctx context.Context, advertiserID string, size, page int) ([]response.Campaign, int64, error) {
	limit, offset := size, size*(page-1)

	_, err := uc.advertiserRepository.GetByID(ctx, advertiserID)
	if err != nil {
		return nil, 0, err
	}

	campaigns, err := uc.campaignRepository.ListWithPagination(ctx, advertiserID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	resp := make([]response.Campaign, len(campaigns))
	for i, campaign := range campaigns {
		// TODO: faster
		var link string
		link, err = uc.s3Client.GetOne(ctx, campaign.ImageID)
		if err != nil {
			return nil, 0, err
		}

		resp[i] = campaign.ToResponse(link)
	}

	count, err := uc.campaignRepository.Count(ctx, advertiserID)
	if err != nil {
		return nil, 0, err
	}

	return resp, count, nil
}
