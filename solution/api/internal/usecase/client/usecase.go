package client

import "api/internal/repository"

// Usecase represents the business logic layer for client operations.
type Usecase struct {
	clientRepository repository.ClientRepository
}

// NewUsecase initializes a new Usecase instance.
func NewUsecase(clientRepository repository.ClientRepository) *Usecase {
	return &Usecase{
		clientRepository: clientRepository,
	}
}
