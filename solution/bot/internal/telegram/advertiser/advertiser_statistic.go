package advertiser

import (
	"bot/pkg/advertising/response"
	"context"

	tele "gopkg.in/telebot.v3"
)

func (h *Handler) advertiserStatistic(c tele.Context) error {
	user, err := h.userRepository.GetByID(context.Background(), c.Sender().ID)
	if err != nil {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "main_menu:back"),
		)
	}

	advertiser, err := h.advertisingClient.GetAdvertiserByID(user.PlatformID)
	if err != nil {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "main_menu:back"),
		)
	}

	statistic, err := h.advertisingClient.GetAdvertiserStats(user.PlatformID)
	if err != nil {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "main_menu:back"),
		)
	}

	return c.Edit(
		h.banners.Advertiser.Caption(h.layout.Text(c, "advertiser_statistic_text", struct {
			response.Statistic
			AdvertiserName string
		}{
			Statistic:      statistic,
			AdvertiserName: advertiser.Name,
		})),
		h.layout.Markup(c, "main_menu:back"),
	)
}
