package advertiser

import (
	"api/internal/entity"
	"api/internal/entity/request"
	"api/internal/entity/response"
	"context"
)

// CreateBulk processes a bulk advertiser creation request.
func (uc *Usecase) CreateBulk(ctx context.Context, req []request.Advertiser) ([]response.Advertiser, error) {
	advertisers := make([]entity.Advertiser, len(req))
	for i, advertiser := range req {
		advertisers[i] = advertiser.ToAdvertiser()
	}

	advertisers, err := uc.advertiserRepository.CreateMany(ctx, advertisers)
	if err != nil {
		return []response.Advertiser{}, err
	}

	resp := make([]response.Advertiser, len(advertisers))
	for i, advertiser := range advertisers {
		resp[i] = advertiser.ToResponse()
	}

	return resp, nil
}
