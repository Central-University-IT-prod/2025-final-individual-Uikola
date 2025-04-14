package server

import (
	"api/internal/server/advertiser"
	"api/internal/server/campaign"
	"api/internal/server/client"
	"api/internal/server/helper"
	"api/internal/server/mlscore"
	"api/internal/server/statistic"
	"api/internal/server/time"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewServer initializes and returns a new HTTP server.
func NewServer(
	clientHandler *client.Handler,
	advertiserHandler *advertiser.Handler,
	mlScoreHandler *mlscore.Handler,
	campaignHandler *campaign.Handler,
	timeHandler *time.Handler,
	statisticHandler *statistic.Handler,
) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)

	addRoutes(router, clientHandler, advertiserHandler, mlScoreHandler, campaignHandler, timeHandler, statisticHandler)

	var handler http.Handler = router

	return handler
}

// registers application routes to the provided router.
func addRoutes(
	router *chi.Mux,
	clientHandler *client.Handler,
	advertiserHandler *advertiser.Handler,
	mlScoreHandler *mlscore.Handler,
	campaignHandler *campaign.Handler,
	timeHandler *time.Handler,
	statisticHandler *statistic.Handler,
) {

	router.Get("/health", health)

	router.Post("/time/advance", helper.MakeHandler(timeHandler.SetDate))
	router.Get("/time/advance", helper.MakeHandler(timeHandler.GetDate))

	router.Route("/clients", func(r chi.Router) {
		r.Post("/bulk", helper.MakeHandler(clientHandler.CreateBulk))
		r.Get("/{clientId}", helper.MakeHandler(clientHandler.GetByID))
	})

	router.Route("/advertisers", func(r chi.Router) {
		r.Post("/bulk", helper.MakeHandler(advertiserHandler.CreateBulk))

		r.Route("/{advertiserId}", func(r chi.Router) {
			r.Get("/", helper.MakeHandler(advertiserHandler.GetByID))
			r.Post("/generate-ad-text", helper.MakeHandler(campaignHandler.GenerateAdText))

			r.Route("/campaigns", func(r chi.Router) {
				r.Post("/", helper.MakeHandler(campaignHandler.Create))
				r.Get("/", helper.MakeHandler(campaignHandler.ListWithPagination))

				r.Route("/{campaignId}", func(r chi.Router) {
					r.Delete("/", helper.MakeHandler(campaignHandler.Delete))
					r.Put("/", helper.MakeHandler(campaignHandler.Update))
					r.Get("/", helper.MakeHandler(campaignHandler.GetByID))
					r.Post("/image", helper.MakeHandler(campaignHandler.UploadImage))
					r.Delete("/image", helper.MakeHandler(campaignHandler.RemoveImage))
				})
			})
		})
	})

	router.Route("/ads", func(r chi.Router) {
		r.Get("/", helper.MakeHandler(campaignHandler.GetAdd))
		r.Post("/{adId}/click", helper.MakeHandler(campaignHandler.ClickAd))
	})

	router.Post("/ml-scores", helper.MakeHandler(mlScoreHandler.Create))

	router.Route("/stats", func(r chi.Router) {
		r.Route("/campaigns", func(r chi.Router) {
			r.Get("/{campaignId}", helper.MakeHandler(statisticHandler.GetCampaign))
			r.Get("/{campaignId}/daily", helper.MakeHandler(statisticHandler.GetCampaignDaily))
		})

		r.Route("/advertisers", func(r chi.Router) {
			r.Get("/{advertiserId}/campaigns", helper.MakeHandler(statisticHandler.GetAdvertiser))
			r.Get("/{advertiserId}/campaigns/daily", helper.MakeHandler(statisticHandler.GetAdvertiserDaily))
		})
	})

	router.Route("/moderation", func(r chi.Router) {
		r.Get("/campaigns", helper.MakeHandler(campaignHandler.ListForModeration))
		r.Post("/moderate/{campaignId}", helper.MakeHandler(campaignHandler.Moderate))
	})
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
