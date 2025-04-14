package e2e

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"gotest.tools/v3/assert"
)

type ClientUpsert struct {
	ClientID string `json:"client_id"`
	Login    string `json:"login"`
	Age      int    `json:"age"`
	Location string `json:"location"`
	Gender   string `json:"gender"`
}

func TestClientLifecycle(t *testing.T) {
	clientID := uuid.New().String()

	t.Run("Create Client", func(t *testing.T) {
		clients := []ClientUpsert{
			{
				ClientID: clientID,
				Login:    "test_user",
				Age:      30,
				Location: "Moscow",
				Gender:   "MALE",
			},
		}

		response := api.POST("/clients/bulk").
			WithJSON(clients).
			Expect().
			Status(http.StatusOK).
			JSON().Array()

		assert.Equal(t, 1, len(response.Raw()), "Ожидался ответ с одним клиентом")
		response.Value(0).Object().Value("client_id").IsEqual(clientID)
	})

	t.Run("Get Client By ID", func(t *testing.T) {
		client := api.GET("/clients/{clientId}", clientID).
			Expect().
			Status(http.StatusOK).
			JSON().Object()

		client.Value("client_id").IsEqual(clientID)
		client.Value("login").IsEqual("test_user")
		client.Value("age").IsEqual(30)
		client.Value("location").IsEqual("Moscow")
		client.Value("gender").IsEqual("MALE")
	})

	t.Run("Update Client", func(t *testing.T) {
		clients := []ClientUpsert{
			{
				ClientID: clientID,
				Login:    "updated_user",
				Age:      35,
				Location: "Saint Petersburg",
				Gender:   "MALE",
			},
		}

		api.POST("/clients/bulk").
			WithJSON(clients).
			Expect().
			Status(http.StatusOK)

		client := api.GET("/clients/{clientId}", clientID).
			Expect().
			Status(http.StatusOK).
			JSON().Object()

		client.Value("login").IsEqual("updated_user")
		client.Value("age").IsEqual(35)
		client.Value("location").IsEqual("Saint Petersburg")
	})

	t.Run("Get Non-Existing Client", func(t *testing.T) {
		api.GET("/clients/{clientId}", uuid.New().String()).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("Invalid Client Data", func(t *testing.T) {
		invalidClient := []map[string]interface{}{
			{
				"client_id": "",
				"login":     "invalid_user",
				"age":       -5,
				"location":  "Unknown",
				"gender":    "UNKNOWN",
			},
		}

		api.POST("/clients/bulk").
			WithJSON(invalidClient).
			Expect().
			Status(http.StatusBadRequest)
	})
}
