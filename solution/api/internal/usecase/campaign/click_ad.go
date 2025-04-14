package campaign

import (
	"api/internal/entity"
	"api/internal/entity/request"
	"api/internal/errorz"
	"context"
)

// ClickAd records a client's click on an advertisement.
func (uc *Usecase) ClickAd(ctx context.Context, req request.ClickAd, adID string) error {
	campaign, err := uc.campaignRepository.GetByID(ctx, adID)
	if err != nil {
		return err
	}

	seen, err := uc.adRepository.HasSeenCampaign(ctx, req.ClientID, campaign.CampaignID)
	if err != nil {
		return err
	}
	if !seen {
		return errorz.ErrCampaignIsNotSeen
	}

	clicked, err := uc.adRepository.HasClickedCampaign(ctx, req.ClientID, campaign.CampaignID)
	if err != nil {
		return err
	}
	if clicked {
		return nil
	}

	currentDay, err := uc.timeRepository.Get(ctx)
	if err != nil {
		return err
	}

	err = uc.adRepository.AddClickedCampaign(ctx, req.ClientID, campaign.CampaignID)
	if err != nil {
		return err
	}

	click := entity.Click{
		AdvertiserID: campaign.AdvertiserID,
		CampaignID:   campaign.CampaignID,
		ClientID:     req.ClientID,
		CostPerClick: campaign.CostPerClick,
		CreatedAt:    currentDay,
	}

	_, err = uc.clickRepository.Create(ctx, click)
	return err
}
