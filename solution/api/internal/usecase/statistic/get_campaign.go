package statistic

import (
	"api/internal/entity/response"
	"context"
)

// GetCampaign retrieves aggregated statistics for a specific campaign.
func (uc *Usecase) GetCampaign(ctx context.Context, campaignID string) (response.GetStatistic, error) {
	_, err := uc.campaignRepository.GetByID(ctx, campaignID)
	if err != nil {
		return response.GetStatistic{}, err
	}

	statistic, err := uc.statisticRepository.GetByCampaignID(ctx, campaignID)
	if err != nil {
		return response.GetStatistic{}, err
	}

	if statistic.ImpressionsCount > 0 {
		statistic.Conversion = (float64(statistic.ClicksCount) / float64(statistic.ImpressionsCount)) * 100
	} else {
		statistic.Conversion = 0
	}
	statistic.SpentTotal = statistic.SpentImpressions.Add(statistic.SpentClicks)

	return statistic.ToGetStatistic(), nil
}
