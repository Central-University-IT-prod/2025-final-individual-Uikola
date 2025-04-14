package campaign

import (
	"api/internal/entity"
	"api/internal/entity/response"
	"api/internal/errorz"
	"api/internal/utils"
	"context"
	"sort"

	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

type campaignWithScore struct {
	Campaign entity.Campaign
	Score    decimal.Decimal
}

// GetAd selects the most suitable advertising campaign for a client.
func (uc *Usecase) GetAd(ctx context.Context, clientID string) (response.GetAd, error) {
	client, err := uc.clientRepository.GetByID(ctx, clientID)
	if err != nil {
		return response.GetAd{}, err
	}

	currentDay, err := uc.timeRepository.Get(ctx)
	if err != nil {
		return response.GetAd{}, err
	}

	campaigns, err := uc.campaignRepository.ListWithTargeting(ctx, client, currentDay)
	if err != nil {
		return response.GetAd{}, err
	}
	if len(campaigns) == 0 {
		return response.GetAd{}, errorz.ErrNoCampaignsFound
	}

	advertisersIDs := make([]string, 0, len(campaigns))
	for _, campaign := range campaigns {
		advertisersIDs = append(advertisersIDs, campaign.AdvertiserID)
	}

	avgMLScore, stddevMLScore, err := uc.mlScoreRepository.AvgSTDDEVWithTargeting(ctx, clientID, advertisersIDs)
	if err != nil {
		return response.GetAd{}, err
	}
	avgCPI, stddevCPI, avgCPC, stddevCPC, err := uc.campaignRepository.AvgSTDDEVCPICPC(ctx, client, currentDay)
	if err != nil {
		return response.GetAd{}, err
	}
	campaignAvgMLScoreMap, campaignSTDDEVScoreMap, err := uc.campaignRepository.GetAvgSTDDEVMLScorePerCampaignWithTargeting(ctx)
	if err != nil {
		return response.GetAd{}, err
	}

	campaignsWithScore := make([]campaignWithScore, len(campaigns))
	for i, campaign := range campaigns {
		var (
			impressionsCount, clicksCount int
			clicked, seen                 bool
		)

		clicked, err = uc.adRepository.HasClickedCampaign(ctx, clientID, campaign.CampaignID)
		if err != nil {
			return response.GetAd{}, err
		}
		if clicked {
			continue
		}

		impressionsCount, err = uc.adRepository.GetImpressionsCount(ctx, campaign.CampaignID)
		if err != nil {
			return response.GetAd{}, err
		}

		clicksCount, err = uc.adRepository.GetClicksCount(ctx, campaign.CampaignID)
		if err != nil {
			return response.GetAd{}, err
		}

		campaignAvgMLScore, ok := campaignAvgMLScoreMap[campaign.CampaignID]
		if !ok {
			continue
		}
		campaignStddevMLScore, ok := campaignSTDDEVScoreMap[campaign.CampaignID]
		if !ok {
			continue
		}

		minMLScore := utils.CalculateMinMLScore(
			decimal.NewFromFloat(campaignAvgMLScore),
			decimal.NewFromFloat(campaignStddevMLScore),
			decimal.NewFromInt(int64(impressionsCount)),
			decimal.NewFromInt(int64(campaign.ImpressionsLimit)),
			decimal.NewFromInt(int64(campaign.EndDate-currentDay)),
			decimal.NewFromInt(int64(campaign.EndDate-campaign.StartDate+1)),
			decimal.NewFromFloat(0.3),
			decimal.NewFromFloat(0.5),
			decimal.NewFromFloat(0.2),
		)

		var mlScore entity.MLScore
		mlScore, err = uc.mlScoreRepository.Get(ctx, clientID, campaign.AdvertiserID)
		if err != nil {
			return response.GetAd{}, err
		}

		if minMLScore.GreaterThan(decimal.NewFromInt(int64(mlScore.Score))) {
			log.Info().Msgf("Skip campaign(%s) because of small ml-score", campaign.AdTitle)
			continue
		}

		seen, err = uc.adRepository.HasSeenCampaign(ctx, clientID, campaign.CampaignID)
		if err != nil {
			continue
		}

		var normCPI decimal.Decimal
		if seen {
			normCPI = decimal.Zero
		} else {
			normCPI = campaign.NormalizeCPI(avgCPI, stddevCPI)
		}
		normCPC := campaign.NormalizeCPC(avgCPC, stddevCPC)

		normMLScore := mlScore.NormalizeScore(decimal.NewFromFloat(avgMLScore), decimal.NewFromFloat(stddevMLScore))

		clickProbability := utils.ProbabilityOfClick(decimal.NewFromInt(int64(mlScore.Score)), decimal.NewFromFloat(avgMLScore), decimal.NewFromFloat(stddevMLScore))

		profit := normCPI.Add(normCPC.Mul(clickProbability))

		impressionsLimitDone := decimal.NewFromFloat(float64(1) - float64(impressionsCount)/float64(campaign.ImpressionsLimit))
		clicksLimitDone := decimal.NewFromFloat(float64(1) - float64(clicksCount)/float64(campaign.ClicksLimit)).Mul(clickProbability)
		if impressionsLimitDone.LessThan(decimal.Zero) {
			impressionsLimitDone = decimal.Zero
		}
		if clicksLimitDone.LessThan(decimal.Zero) {
			clicksLimitDone = decimal.Zero
		}

		limitDone := impressionsLimitDone.Add(clicksLimitDone)

		score := profit.Add(normMLScore.Mul(decimal.NewFromFloat(0.5))).Add(limitDone.Mul(decimal.NewFromFloat(0.07)))

		campaignsWithScore[i] = campaignWithScore{
			Campaign: campaign,
			Score:    score,
		}
	}

	sort.Slice(campaignsWithScore, func(i, j int) bool {
		return campaignsWithScore[j].Score.LessThan(campaignsWithScore[i].Score)
	})

	var allowed bool
	var bestCampaign entity.Campaign
	for _, bestCampaignWithScore := range campaignsWithScore {
		var seen bool
		seen, err = uc.adRepository.HasSeenCampaign(ctx, clientID, bestCampaignWithScore.Campaign.CampaignID)
		if err != nil {
			return response.GetAd{}, err
		}
		if seen {
			bestCampaign = bestCampaignWithScore.Campaign
			break
		}

		allowed, err = uc.adRepository.TryIncrementImpressions(ctx, bestCampaignWithScore.Campaign.CampaignID, bestCampaignWithScore.Campaign.ImpressionsLimit)
		if err != nil {
			return response.GetAd{}, err
		}
		if allowed {
			bestCampaign = bestCampaignWithScore.Campaign
			break
		}
	}

	if bestCampaign.CampaignID == "" {
		return response.GetAd{}, errorz.ErrNoCampaignsFound
	}

	err = uc.adRepository.AddSeenCampaign(ctx, clientID, bestCampaign.CampaignID)
	if err != nil {
		return response.GetAd{}, err
	}

	_, err = uc.impressionRepository.Create(ctx, entity.Impression{
		AdvertiserID:      bestCampaign.AdvertiserID,
		CampaignID:        bestCampaign.CampaignID,
		ClientID:          clientID,
		CostPerImpression: bestCampaign.CostPerImpression,
		CreatedAt:         currentDay,
	})
	if err != nil {
		return response.GetAd{}, err
	}

	link, err := uc.s3Client.GetOne(ctx, bestCampaign.ImageID)
	if err != nil {
		return response.GetAd{}, err
	}

	resp := response.GetAd{
		AdID:         bestCampaign.CampaignID,
		AdTitle:      bestCampaign.AdTitle,
		AdText:       bestCampaign.AdText,
		AdvertiserID: bestCampaign.AdvertiserID,
		ImageURL:     link,
	}

	return resp, nil
}
