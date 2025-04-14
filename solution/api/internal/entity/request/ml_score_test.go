package request_test

import (
	"api/internal/entity/request"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMLScore_Validate(t *testing.T) {
	validMLScore := request.MLScore{
		ClientID:     "123e4567-e89b-12d3-a456-426614174000",
		AdvertiserID: "223e4567-e89b-12d3-a456-426614174000",
		Score:        100,
	}
	assert.NoError(t, validMLScore.Validate())

	invalidMLScore := request.MLScore{ClientID: "invalid-uuid"}
	assert.Error(t, invalidMLScore.Validate())
}

func TestMLScore_ToMLScore(t *testing.T) {
	r := request.MLScore{
		ClientID:     "123e4567-e89b-12d3-a456-426614174000",
		AdvertiserID: "223e4567-e89b-12d3-a456-426614174000",
		Score:        100,
	}

	mlScore := r.ToMLScore()
	assert.Equal(t, "123e4567-e89b-12d3-a456-426614174000", mlScore.ClientID)
	assert.Equal(t, "223e4567-e89b-12d3-a456-426614174000", mlScore.AdvertiserID)
	assert.Equal(t, 100, mlScore.Score)
}
