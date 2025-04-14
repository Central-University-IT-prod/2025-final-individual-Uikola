package campaign

import (
	"api/internal/entity/response"
	"context"
)

// ListForModerationWithPagination retrieves a list of campaigns for moderation with pagination.
func (uc *Usecase) ListForModerationWithPagination(ctx context.Context, size, page int) ([]response.Campaign, int64, error) {
	limit, offset := size, size*(page-1)

	campaigns, err := uc.campaignRepository.ListForModerationWithPagination(ctx, limit, offset)
	if err != nil {
		return []response.Campaign{}, 0, err
	}

	resp := make([]response.Campaign, len(campaigns))
	for i, campaign := range campaigns {
		var link string
		link, err = uc.s3Client.GetOne(ctx, campaign.ImageID)
		if err != nil {
			return []response.Campaign{}, 0, err
		}

		resp[i] = campaign.ToResponse(link)
	}

	count, err := uc.campaignRepository.CountForModeration(ctx)
	if err != nil {
		return []response.Campaign{}, 0, err
	}

	return resp, count, nil
}
