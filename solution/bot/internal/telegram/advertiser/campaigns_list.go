package advertiser

import (
	"bot/pkg/advertising/response"
	"context"

	"github.com/rs/zerolog/log"
	tele "gopkg.in/telebot.v3"
)

func (h *Handler) campaignsList(c tele.Context) error {
	const campaignsOnPage = 5

	var (
		p              int
		prevPage       int
		nextPage       int
		err            error
		campaignsCount int
		campaigns      []response.Campaign
		rows           []tele.Row
		menuRow        tele.Row
	)

	_, p, err = parseCallback(c.Callback().Data)
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
			h.layout.Markup(c, "main_menu:back"),
		)
	}

	campaigns, campaignsCount, err = h.advertisingClient.ListCampaigns(user.PlatformID, campaignsOnPage, p)
	if err != nil {
		log.Error().Msgf("(user: %d) error while get campaigns: %v", c.Sender().ID, err)
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "main_menu:back"),
		)
	}

	currentDate := 2

	markup := c.Bot().NewMarkup()
	for _, campaign := range campaigns {
		rows = append(rows, markup.Row(*h.layout.Button(c, "campaigns:campaign", struct {
			ID     string
			Page   int
			Title  string
			IsOver bool
		}{
			ID:     campaign.CampaignID,
			Page:   p,
			Title:  campaign.AdTitle,
			IsOver: campaign.IsOver(currentDate),
		})))
	}
	pagesCount := ((campaignsCount - 1) / campaignsOnPage) + 1
	if p == 1 {
		prevPage = pagesCount
	} else {
		prevPage = p - 1
	}

	if p >= pagesCount {
		nextPage = 1
	} else {
		nextPage = p + 1
	}

	menuRow = append(menuRow,
		*h.layout.Button(c, "advertiser:campaigns:prev_page", struct {
			Page int
		}{
			Page: prevPage,
		}),
		*h.layout.Button(c, "core:page_counter", struct {
			Page       int
			PagesCount int
		}{
			Page:       p,
			PagesCount: pagesCount,
		}),
		*h.layout.Button(c, "advertiser:campaigns:next_page", struct {
			ID   string
			Page int
		}{
			Page: nextPage,
		}),
	)

	rows = append(
		rows,
		menuRow,
		markup.Row(*h.layout.Button(c, "main_menu:back")),
	)

	markup.Inline(rows...)

	return c.Edit(
		h.banners.Advertiser.Caption(h.layout.Text(c, "campaigns_list")),
		markup,
	)
}
