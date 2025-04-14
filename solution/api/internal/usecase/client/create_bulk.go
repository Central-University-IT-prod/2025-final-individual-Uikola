package client

import (
	"api/internal/entity"
	"api/internal/entity/request"
	"api/internal/entity/response"
	"context"
)

// CreateBulk processes a bulk client creation request.
func (uc *Usecase) CreateBulk(ctx context.Context, req []request.Client) ([]response.Client, error) {
	clients := make([]entity.Client, len(req))
	for i, client := range req {
		clients[i] = client.ToClient()
	}

	clients, err := uc.clientRepository.CreateMany(ctx, clients)
	if err != nil {
		return []response.Client{}, err
	}

	resp := make([]response.Client, len(clients))
	for i, client := range clients {
		resp[i] = client.ToResponse()
	}

	return resp, nil
}
