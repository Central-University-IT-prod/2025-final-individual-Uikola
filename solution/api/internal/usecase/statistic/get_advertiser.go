package statistic

import (
	"api/internal/entity/response"
	"context"
)

// GetAdvertiser retrieves the aggregated statistics for a specific advertiser.
func (uc *Usecase) GetAdvertiser(ctx context.Context, advertiserID string) (response.GetStatistic, error) {
	_, err := uc.advertiserRepository.GetByID(ctx, advertiserID)
	if err != nil {
		return response.GetStatistic{}, err
	}

	statistic, err := uc.statisticRepository.GetByAdvertiserID(ctx, advertiserID)
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
