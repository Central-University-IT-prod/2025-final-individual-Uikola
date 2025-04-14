package campaign

import (
	"api/internal/repository"
	"api/pkg/ai"
	"api/pkg/s3"
)

// Usecase represents the business logic layer for campaign operations.
type Usecase struct {
	campaignRepository   repository.CampaignRepository
	timeRepository       repository.TimeRepository
	clientRepository     repository.ClientRepository
	impressionRepository repository.ImpressionRepository
	clickRepository      repository.ClickRepository
	mlScoreRepository    repository.MLScoreRepository
	advertiserRepository repository.AdvertiserRepository
	statisticRepository  repository.StatisticRepository
	adRepository         repository.AdRepository

	s3Client s3.Client
	aiClient ai.Client
}

// NewUsecase initializes a new Usecase instance.
func NewUsecase(
	campaignRepository repository.CampaignRepository,
	timeRepository repository.TimeRepository,
	clientRepository repository.ClientRepository,
	impressionRepository repository.ImpressionRepository,
	clickRepository repository.ClickRepository,
	mlScoreRepository repository.MLScoreRepository,
	advertiserRepository repository.AdvertiserRepository,
	statisticRepository repository.StatisticRepository,
	adRepository repository.AdRepository,
	s3Client s3.Client,
	aiClient ai.Client,
) *Usecase {
	return &Usecase{
		campaignRepository:   campaignRepository,
		timeRepository:       timeRepository,
		clientRepository:     clientRepository,
		impressionRepository: impressionRepository,
		clickRepository:      clickRepository,
		mlScoreRepository:    mlScoreRepository,
		advertiserRepository: advertiserRepository,
		statisticRepository:  statisticRepository,
		adRepository:         adRepository,

		s3Client: s3Client,
		aiClient: aiClient,
	}
}

func toString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
