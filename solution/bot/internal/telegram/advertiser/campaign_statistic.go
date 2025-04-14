package advertiser

import tele "gopkg.in/telebot.v3"

func (h *Handler) campaignStatistic(c tele.Context) error {
	campaignID, p, err := parseCallback(c.Callback().Data)
	if err != nil {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "main_menu:back"),
		)
	}

	statistic, err := h.advertisingClient.GetCampaignStats(campaignID)
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
		h.banners.Advertiser.Caption(h.layout.Text(c, "campaign_statistic_text", statistic)),
		h.layout.Markup(c, "campaigns:campaign:back", struct {
			ID   string
			Page int
		}{
			ID:   campaignID,
			Page: p,
		}),
	)
}
