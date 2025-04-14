package campaign

import (
	"api/internal/entity"
	"api/internal/entity/request"
	"api/internal/entity/response"
	"api/internal/errorz"
	"context"

	"github.com/shopspring/decimal"
)

// Update processes a campaign update request.
func (uc *Usecase) Update(ctx context.Context, advertiserID, campaignID string, req request.UpdateCampaign) (response.Campaign, error) {
	_, err := uc.advertiserRepository.GetByID(ctx, advertiserID)
	if err != nil {
		return response.Campaign{}, err
	}

	campaign, err := uc.campaignRepository.Get(ctx, advertiserID, campaignID)
	if err != nil {
		return response.Campaign{}, err
	}

	currentDay, err := uc.timeRepository.Get(ctx)
	if err != nil {
		return response.Campaign{}, err
	}

	if campaign.IsOver(currentDay) {
		return response.Campaign{}, errorz.ErrCampaignIsOver
	}

	updatedCampaign, err := uc.update(campaign, req, currentDay)
	if err != nil {
		return response.Campaign{}, err
	}

	err = uc.campaignRepository.Update(ctx, updatedCampaign)
	if err != nil {
		return response.Campaign{}, err
	}

	campaign, err = uc.campaignRepository.Get(ctx, advertiserID, campaignID)
	if err != nil {
		return response.Campaign{}, err
	}

	link, err := uc.s3Client.GetOne(ctx, campaign.ImageID)
	if err != nil {
		return response.Campaign{}, err
	}

	return campaign.ToResponse(link), nil
}

func (uc *Usecase) update(oldCampaign entity.Campaign, req request.UpdateCampaign, currentDay int) (entity.Campaign, error) {
	if oldCampaign.IsActive(currentDay) {
		if req.ImpressionsLimit != oldCampaign.ImpressionsLimit {
			return entity.Campaign{}, errorz.ErrCampaignIsActive
		}
		if req.ClicksLimit != oldCampaign.ClicksLimit {
			return entity.Campaign{}, errorz.ErrCampaignIsActive
		}
		if req.StartDate != oldCampaign.StartDate {
			return entity.Campaign{}, errorz.ErrCampaignIsActive
		}
		if req.EndDate != oldCampaign.EndDate {
			return entity.Campaign{}, errorz.ErrCampaignIsActive
		}
	} else {
		if req.StartDate < oldCampaign.StartDate {
			return entity.Campaign{}, errorz.ErrInvalidDate
		}
		if req.EndDate <= req.StartDate {
			return entity.Campaign{}, errorz.ErrInvalidDate
		}
	}

	campaign := entity.Campaign{
		CampaignID:        oldCampaign.CampaignID,
		AdvertiserID:      oldCampaign.AdvertiserID,
		ImpressionsLimit:  req.ImpressionsLimit,
		ClicksLimit:       req.ClicksLimit,
		CostPerImpression: decimal.NewFromFloat(req.CostPerImpression),
		CostPerClick:      decimal.NewFromFloat(req.CostPerClick),
		AdTitle:           toString(req.AdTitle),
		AdText:            toString(req.AdText),
		StartDate:         req.StartDate,
		EndDate:           req.EndDate,
		ImageID:           oldCampaign.ImageID,
		PassedModeration:  false,
		CreatedAt:         oldCampaign.CreatedAt,
		Targeting: entity.Targeting{
			TargetingID: oldCampaign.Targeting.TargetingID,
			CampaignID:  oldCampaign.CampaignID,
		},
	}

	if req.Targeting != nil {
		if req.Targeting.Gender != nil {
			campaign.Targeting.Gender = req.Targeting.Gender
		}
		if req.Targeting.AgeFrom != nil {
			campaign.Targeting.AgeFrom = req.Targeting.AgeFrom
		}
		if req.Targeting.AgeTo != nil {
			campaign.Targeting.AgeTo = req.Targeting.AgeTo
		}
		if req.Targeting.Location != nil {
			campaign.Targeting.Location = req.Targeting.Location
		}
	}

	return campaign, nil
}
