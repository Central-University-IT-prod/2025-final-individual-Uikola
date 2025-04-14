package statistic

import "api/internal/repository"

// Usecase represents the business logic layer for statistic operations.
type Usecase struct {
	statisticRepository  repository.StatisticRepository
	advertiserRepository repository.AdvertiserRepository
	campaignRepository   repository.CampaignRepository
}

// NewUsecase initializes a new Usecase instance.
func NewUsecase(statisticRepository repository.StatisticRepository, advertiserRepository repository.AdvertiserRepository, campaignRepository repository.CampaignRepository) *Usecase {
	return &Usecase{
		statisticRepository:  statisticRepository,
		advertiserRepository: advertiserRepository,
		campaignRepository:   campaignRepository,
	}
}
