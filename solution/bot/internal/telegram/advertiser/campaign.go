package advertiser

import (
	"bot/pkg/advertising/response"
	"context"

	"github.com/rs/zerolog/log"
	tele "gopkg.in/telebot.v3"
)

func (h *Handler) campaign(c tele.Context) error {
	campaignID, p, err := parseCallback(c.Callback().Data)
	if err != nil {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "main_menu:back"),
		)
	}

	user, err := h.userRepository.GetByID(context.Background(), c.Sender().ID)
	if err != nil {
		log.Error().Msgf("(user: %d) error while get user: %v", c.Sender().ID, err)
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "advertiser:campaign:back", struct {
				Page int
			}{
				Page: p,
			}),
		)
	}

	campaign, err := h.advertisingClient.GetCampaign(user.PlatformID, campaignID)
	if err != nil {
		log.Error().Msgf("(user: %d) error while get campaign: %v", c.Sender().ID, err)
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),

			h.layout.Markup(c, "advertiser:campaign:back", struct {
				Page int
			}{
				Page: p,
			}))
	}

	currentDate, err := h.advertisingClient.GetCurrentDate()
	if err != nil {
		log.Error().Msgf("(user: %d) error while get current date: %v", c.Sender().ID, err)
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),

			h.layout.Markup(c, "advertiser:campaign:back", struct {
				Page int
			}{
				Page: p,
			}))
	}

	return c.Edit(
		h.banners.Advertiser.Caption(h.layout.Text(c, "campaign", struct {
			response.Campaign
			CurrentDate int
		}{
			Campaign:    campaign,
			CurrentDate: currentDate,
		})),
		h.layout.Markup(c, "advertiser:campaign:menu", struct {
			ID   string
			Page int
		}{
			ID:   campaignID,
			Page: p,
		}),
	)
}
