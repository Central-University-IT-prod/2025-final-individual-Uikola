package e2e

import (
	"net/http"
	"testing"
)

type AdvanceRequest struct {
	CurrentDate int `json:"current_date"`
}

type AdvanceResponse struct {
	CurrentDate int `json:"current_date"`
}

func TestTimeLifecycle(t *testing.T) {
	t.Run("Set Current Date", func(t *testing.T) {
		request := AdvanceRequest{
			CurrentDate: 1,
		}

		response := api.POST("/time/advance").
			WithJSON(request).
			Expect().
			Status(http.StatusOK).
			JSON().Object()

		response.Value("current_date").IsEqual(request.CurrentDate)
	})

	t.Run("Get Current Date", func(t *testing.T) {
		response := api.GET("/time/advance").
			Expect().
			Status(http.StatusOK).
			JSON().Object()

		response.Value("current_date").IsNumber().IsEqual(1)
	})

	t.Run("Set Invalid Date", func(t *testing.T) {
		request := AdvanceRequest{
			CurrentDate: -5,
		}

		api.POST("/time/advance").
			WithJSON(request).
			Expect().
			Status(http.StatusBadRequest)
	})
}
