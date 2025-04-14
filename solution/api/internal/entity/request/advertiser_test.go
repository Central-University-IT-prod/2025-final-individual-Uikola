package request_test

import (
	"api/internal/entity"
	"api/internal/entity/request"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAdvertiser_Validate(t *testing.T) {
	tests := []struct {
		name    string
		input   request.Advertiser
		expects bool
	}{
		{
			name: "Valid Advertiser",
			input: request.Advertiser{
				AdvertiserID: "550e8400-e29b-41d4-a716-446655440000",
				Name:         "Advertiser One",
			},
			expects: true,
		},
		{
			name: "Missing AdvertiserID",
			input: request.Advertiser{
				Name: "Advertiser Two",
			},
			expects: false,
		},
		{
			name: "Invalid UUID",
			input: request.Advertiser{
				AdvertiserID: "invalid-uuid",
				Name:         "Advertiser Three",
			},
			expects: false,
		},
		{
			name: "Missing Name",
			input: request.Advertiser{
				AdvertiserID: "550e8400-e29b-41d4-a716-446655440000",
			},
			expects: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if tt.expects {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestAdvertiser_ToAdvertiser(t *testing.T) {
	input := request.Advertiser{
		AdvertiserID: "550e8400-e29b-41d4-a716-446655440000",
		Name:         "Test Advertiser",
	}

	expected := entity.Advertiser{
		AdvertiserID: "550e8400-e29b-41d4-a716-446655440000",
		Name:         "Test Advertiser",
	}

	result := input.ToAdvertiser()

	assert.Equal(t, expected, result)
}
