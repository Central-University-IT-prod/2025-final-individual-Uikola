package statistic

import (
	"api/internal/entity/response"
	"context"
)

// GetCampaignDaily retrieves daily statistics for a specific campaign.
func (uc *Usecase) GetCampaignDaily(ctx context.Context, campaignID string) ([]response.DailyStatistic, error) {
	_, err := uc.campaignRepository.GetByID(ctx, campaignID)
	if err != nil {
		return nil, err
	}

	statistics, err := uc.statisticRepository.GetDailyByCampaignID(ctx, campaignID)
	if err != nil {
		return []response.DailyStatistic{}, err
	}

	resp := make([]response.DailyStatistic, len(statistics))

	for i, statistic := range statistics {
		if statistic.ImpressionsCount > 0 {
			statistic.Conversion = (float64(statistic.ClicksCount) / float64(statistic.ImpressionsCount)) * 100
		} else {
			statistic.Conversion = 0
		}
		statistic.SpentTotal = statistic.SpentImpressions.Add(statistic.SpentClicks)
		resp[i] = statistic.ToDailyStatistic()
	}
	return resp, nil
}
