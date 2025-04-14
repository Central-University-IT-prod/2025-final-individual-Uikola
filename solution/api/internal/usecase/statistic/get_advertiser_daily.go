package statistic

import (
	"api/internal/entity/response"
	"context"
)

// GetAdvertiserDaily retrieves daily statistics for a specific advertiser.
func (uc *Usecase) GetAdvertiserDaily(ctx context.Context, advertiserID string) ([]response.DailyStatistic, error) {
	_, err := uc.advertiserRepository.GetByID(ctx, advertiserID)
	if err != nil {
		return nil, err
	}

	statistics, err := uc.statisticRepository.GetDailyByAdvertiserID(ctx, advertiserID)
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
