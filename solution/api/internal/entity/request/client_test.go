package request_test

import (
	"api/internal/entity"
	"api/internal/entity/request"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient_Validate(t *testing.T) {
	validClient := request.Client{
		ClientID: "123e4567-e89b-12d3-a456-426614174000",
		Login:    "user123",
		Age:      30,
		Location: "New York",
		Gender:   entity.ClientGenderMale,
	}
	assert.NoError(t, validClient.Validate())

	invalidClient := request.Client{Login: "user123"}
	assert.Error(t, invalidClient.Validate())
}

func TestClient_ToClient(t *testing.T) {
	r := request.Client{
		ClientID: "123e4567-e89b-12d3-a456-426614174000",
		Login:    "user123",
		Age:      30,
		Location: "New York",
		Gender:   entity.ClientGenderMale,
	}

	client := r.ToClient()
	assert.Equal(t, "123e4567-e89b-12d3-a456-426614174000", client.ClientID)
	assert.Equal(t, "user123", client.Login)
	assert.Equal(t, 30, client.Age)
	assert.Equal(t, "New York", client.Location)
	assert.Equal(t, entity.ClientGenderMale, client.Gender)
}
