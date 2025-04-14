package advertiser

import (
	"context"

	tele "gopkg.in/telebot.v3"
)

func (h *Handler) deleteCampaign(c tele.Context) error {
	campaignID, p, err := parseCallback(c.Callback().Data)
	if err != nil {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "main_menu:back"),
		)
	}

	user, err := h.userRepository.GetByID(context.Background(), c.Sender().ID)
	if err != nil {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "campaigns:campaign:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	campaign, err := h.advertisingClient.GetCampaign(user.PlatformID, campaignID)
	if err != nil {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "campaigns:campaign:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	return c.Edit(
		h.banners.Advertiser.Caption(h.layout.Text(c, "campaign_delete_confirmation", campaign)),
		h.layout.Markup(c, "campaign:delete:confirmation", struct {
			ID   string
			Page int
		}{
			ID:   campaignID,
			Page: p,
		}),
	)
}

func (h *Handler) confirmCampaignDeletion(c tele.Context) error {
	campaignID, p, err := parseCallback(c.Callback().Data)
	if err != nil {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "main_menu:back"),
		)
	}

	user, err := h.userRepository.GetByID(context.Background(), c.Sender().ID)
	if err != nil {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "campaigns:campaign:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	campaign, err := h.advertisingClient.GetCampaign(user.PlatformID, campaignID)
	if err != nil {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "campaigns:campaign:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	err = h.advertisingClient.DeleteCampaign(user.PlatformID, campaignID)
	if err != nil {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "campaigns:campaign:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	return c.Edit(
		h.banners.Advertiser.Caption(h.layout.Text(c, "campaign_deleted_successfully", campaign)),
		h.layout.Markup(c, "main_menu:back"),
	)
}
