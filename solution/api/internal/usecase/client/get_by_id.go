package client

import (
	"api/internal/entity/response"
	"context"
)

// GetByID retrieves a client by their client ID.
func (uc *Usecase) GetByID(ctx context.Context, clientID string) (response.Client, error) {
	client, err := uc.clientRepository.GetByID(ctx, clientID)
	if err != nil {
		return response.Client{}, err
	}

	return client.ToResponse(), nil
}
