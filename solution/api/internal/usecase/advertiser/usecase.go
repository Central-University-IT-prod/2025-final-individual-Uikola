package advertiser

import "api/internal/repository"

// Usecase represents the business logic layer for advertiser operations.
type Usecase struct {
	advertiserRepository repository.AdvertiserRepository
}

// NewUsecase initializes a new Usecase instance.
func NewUsecase(advertiserRepository repository.AdvertiserRepository) *Usecase {
	return &Usecase{
		advertiserRepository: advertiserRepository,
	}
}
