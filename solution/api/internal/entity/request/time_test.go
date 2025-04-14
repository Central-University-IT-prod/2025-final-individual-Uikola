package request_test

import (
	"api/internal/entity/request"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAdvance_Validate(t *testing.T) {
	validAdvance := request.Advance{
		CurrentDate: 5,
	}
	assert.NoError(t, validAdvance.Validate())

	invalidAdvance := request.Advance{
		CurrentDate: -1,
	}
	assert.Error(t, invalidAdvance.Validate())
}
