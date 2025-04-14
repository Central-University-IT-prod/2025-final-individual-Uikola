package e2e

import (
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gavv/httpexpect/v2"
	"github.com/joho/godotenv"
)

var api *httpexpect.Expect
var baseURL string

func TestMain(m *testing.M) {
	if err := godotenv.Load("e2e.env"); err != nil {
		log.Fatalf("Error loading e2e.env: %v", err)
	}

	baseURL = os.Getenv("API_URL")
	if baseURL == "" {
		log.Fatal("API_URL is not set in e2e.env")
	}

	api = httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  baseURL,
		Reporter: httpexpect.NewPanicReporter(),
		Client:   &http.Client{Timeout: 30 * time.Second},
	})

	waitForAPI()

	code := m.Run()
	os.Exit(code)
}

func waitForAPI() {
	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			log.Fatalf("Timeout: API not ready at %s", baseURL)
		case <-ticker.C:
			resp := api.GET("/health").Expect().Raw()
			if resp.StatusCode == 200 {
				return
			}
		}
	}
}
