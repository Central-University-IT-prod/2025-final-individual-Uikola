package advertiser

import (
	"bot/internal/telegram/helper/validator"
	"bot/pkg/advertising/errorz"
	"bot/pkg/advertising/request"
	"context"
	"errors"
	"strconv"

	"github.com/nlypage/intele"
	"github.com/nlypage/intele/collector"
	"github.com/rs/zerolog/log"
	tele "gopkg.in/telebot.v3"
)

func (h *Handler) updateCampaign(c tele.Context) error {
	campaignID, p, err := parseCallback(c.Callback().Data)
	if err != nil {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "main_menu:back"),
		)
	}

	currentDate, err := h.advertisingClient.GetCurrentDate()
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
		h.banners.Advertiser.Caption(h.layout.Text(c, "campaign_update_text", struct {
			CurrentDate int
		}{
			CurrentDate: currentDate,
		})),
		h.layout.Markup(c, "campaign:update:menu", struct {
			ID   string
			Page int
		}{
			ID:   campaignID,
			Page: p,
		}),
	)
}

func (h *Handler) updateImpressionsLimit(c tele.Context) error {
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
			h.layout.Markup(c, "campaign:update:back", struct {
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
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	currentDate, err := h.advertisingClient.GetCurrentDate()
	if err != nil {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	if campaign.IsActive(currentDate) {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "campaign_is_active")),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}
	if campaign.IsOver(currentDate) {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "campaign_is_over")),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	inputCollector := collector.New()
	_ = c.Edit(
		h.banners.Advertiser.Caption(h.layout.Text(c, "input_impressions_limit")),
		h.layout.Markup(c, "campaign:update:back", struct {
			ID   string
			Page int
		}{
			ID:   campaignID,
			Page: p,
		}),
	)
	inputCollector.Collect(c.Message())

	var (
		impressionsLimit int
		done             bool
		resp             intele.Response
	)
	for {
		resp, err = h.inputManager.Get(context.Background(), c.Sender().ID, 0)
		if resp.Message != nil {
			inputCollector.Collect(resp.Message)
		}
		switch {
		case resp.Canceled:
			_ = inputCollector.Clear(c, collector.ClearOptions{IgnoreErrors: true, ExcludeLast: true})
			return nil
		case err != nil:
			log.Error().Msgf("(user: %d) error while input impressions limit: %v", c.Sender().ID, err)
			_ = inputCollector.Send(c,
				h.banners.Advertiser.Caption(h.layout.Text(c, "input_error", h.layout.Text(c, "input_impressions_limit"))),
				h.layout.Markup(c, "campaign:update:back", struct {
					ID   string
					Page int
				}{
					ID:   campaignID,
					Page: p,
				}),
			)
		case resp.Message == nil:
			log.Error().Msgf("(user: %d) error while input impressions limit: %v", c.Sender().ID, err)
			_ = inputCollector.Send(c,
				h.banners.Advertiser.Caption(h.layout.Text(c, "input_error", h.layout.Text(c, "input_impressions_limit"))),
				h.layout.Markup(c, "campaign:update:back", struct {
					ID   string
					Page int
				}{
					ID:   campaignID,
					Page: p,
				}),
			)
		case !validator.ImpressionsLimit(resp.Message.Text, nil):
			_ = inputCollector.Send(c,
				h.banners.Advertiser.Caption(h.layout.Text(c, "invalid_impressions_limit")),
				h.layout.Markup(c, "campaign:update:back", struct {
					ID   string
					Page int
				}{
					ID:   campaignID,
					Page: p,
				}),
			)
		case validator.ImpressionsLimit(resp.Message.Text, nil):
			impressionsLimit, _ = strconv.Atoi(resp.Message.Text)
			_ = inputCollector.Clear(c, collector.ClearOptions{IgnoreErrors: true})
			done = true
		}
		if done {
			break
		}
	}

	updateRequestBody := request.UpdateCampaign{
		ImpressionsLimit:  impressionsLimit,
		ClicksLimit:       campaign.ClicksLimit,
		CostPerImpression: campaign.CostPerImpression,
		CostPerClick:      campaign.CostPerClick,
		AdTitle:           &campaign.AdTitle,
		AdText:            &campaign.AdText,
		StartDate:         campaign.StartDate,
		EndDate:           campaign.EndDate,
		Targeting: &request.Targeting{
			Gender:   campaign.Targeting.Gender,
			AgeFrom:  campaign.Targeting.AgeFrom,
			AgeTo:    campaign.Targeting.AgeTo,
			Location: campaign.Targeting.Location,
		},
	}

	updatedCampaign, err := h.advertisingClient.UpdateCampaign(updateRequestBody, user.PlatformID, campaignID)
	switch {
	case errors.Is(err, errorz.ErrInvalidData) || errors.Is(err, errorz.ErrCampaignNotFound):
		return c.Send(
			h.banners.Advertiser.Caption(h.layout.Text(c, "invalid_impressions_limit")),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	case err != nil:
		return c.Send(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	return c.Send(
		h.banners.Advertiser.Caption(h.layout.Text(c, "campaign_updated_successfully", struct {
			AdTitle string
		}{
			AdTitle: updatedCampaign.AdTitle,
		})),
		h.layout.Markup(c, "campaign:update:back", struct {
			ID   string
			Page int
		}{
			ID:   campaignID,
			Page: p,
		}),
	)
}

func (h *Handler) updateClicksLimit(c tele.Context) error {
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
			h.layout.Markup(c, "campaign:update:back", struct {
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
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	currentDate, err := h.advertisingClient.GetCurrentDate()
	if err != nil {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	if campaign.IsActive(currentDate) {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "campaign_is_active")),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}
	if campaign.IsOver(currentDate) {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "campaign_is_over")),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	inputCollector := collector.New()
	_ = c.Edit(
		h.banners.Advertiser.Caption(h.layout.Text(c, "input_clicks_limit")),
		h.layout.Markup(c, "campaign:update:back", struct {
			ID   string
			Page int
		}{
			ID:   campaignID,
			Page: p,
		}),
	)
	inputCollector.Collect(c.Message())

	var (
		clicksLimit int
		done        bool
		resp        intele.Response
	)
	for {
		resp, err = h.inputManager.Get(context.Background(), c.Sender().ID, 0)
		if resp.Message != nil {
			inputCollector.Collect(resp.Message)
		}
		switch {
		case resp.Canceled:
			_ = inputCollector.Clear(c, collector.ClearOptions{IgnoreErrors: true, ExcludeLast: true})
			return nil
		case err != nil:
			log.Error().Msgf("(user: %d) error while input clicks limit: %v", c.Sender().ID, err)
			_ = inputCollector.Send(c,
				h.banners.Advertiser.Caption(h.layout.Text(c, "input_error", h.layout.Text(c, "input_clicks_limit"))),
				h.layout.Markup(c, "campaign:update:back", struct {
					ID   string
					Page int
				}{
					ID:   campaignID,
					Page: p,
				}),
			)
		case resp.Message == nil:
			log.Error().Msgf("(user: %d) error while input clicks limit: %v", c.Sender().ID, err)
			_ = inputCollector.Send(c,
				h.banners.Advertiser.Caption(h.layout.Text(c, "input_error", h.layout.Text(c, "input_clicks_limit"))),
				h.layout.Markup(c, "campaign:update:back", struct {
					ID   string
					Page int
				}{
					ID:   campaignID,
					Page: p,
				}),
			)
		case !validator.ClicksLimit(resp.Message.Text, nil):
			_ = inputCollector.Send(c,
				h.banners.Advertiser.Caption(h.layout.Text(c, "invalid_clicks_limit")),
				h.layout.Markup(c, "campaign:update:back", struct {
					ID   string
					Page int
				}{
					ID:   campaignID,
					Page: p,
				}),
			)
		case validator.ImpressionsLimit(resp.Message.Text, nil):
			clicksLimit, _ = strconv.Atoi(resp.Message.Text)
			_ = inputCollector.Clear(c, collector.ClearOptions{IgnoreErrors: true})
			done = true
		}
		if done {
			break
		}
	}

	updateRequestBody := request.UpdateCampaign{
		ImpressionsLimit:  campaign.ImpressionsLimit,
		ClicksLimit:       clicksLimit,
		CostPerImpression: campaign.CostPerImpression,
		CostPerClick:      campaign.CostPerClick,
		AdTitle:           &campaign.AdTitle,
		AdText:            &campaign.AdText,
		StartDate:         campaign.StartDate,
		EndDate:           campaign.EndDate,
		Targeting: &request.Targeting{
			Gender:   campaign.Targeting.Gender,
			AgeFrom:  campaign.Targeting.AgeFrom,
			AgeTo:    campaign.Targeting.AgeTo,
			Location: campaign.Targeting.Location,
		},
	}

	updatedCampaign, err := h.advertisingClient.UpdateCampaign(updateRequestBody, user.PlatformID, campaignID)
	switch {
	case errors.Is(err, errorz.ErrInvalidData) || errors.Is(err, errorz.ErrCampaignNotFound):
		return c.Send(
			h.banners.Advertiser.Caption(h.layout.Text(c, "invalid_clicks_limit")),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	case err != nil:
		return c.Send(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	return c.Send(
		h.banners.Advertiser.Caption(h.layout.Text(c, "campaign_updated_successfully", struct {
			AdTitle string
		}{
			AdTitle: updatedCampaign.AdTitle,
		})),
		h.layout.Markup(c, "campaign:update:back", struct {
			ID   string
			Page int
		}{
			ID:   campaignID,
			Page: p,
		}),
	)
}

func (h *Handler) updateCostPerImpression(c tele.Context) error {
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
			h.layout.Markup(c, "campaign:update:back", struct {
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
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)

	}

	currentDate, err := h.advertisingClient.GetCurrentDate()
	if err != nil {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	if campaign.IsOver(currentDate) {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "campaign_is_over")),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	inputCollector := collector.New()
	_ = c.Edit(
		h.banners.Advertiser.Caption(h.layout.Text(c, "input_cost_per_impression")),
		h.layout.Markup(c, "campaign:update:back", struct {
			ID   string
			Page int
		}{
			ID:   campaignID,
			Page: p,
		}),
	)
	inputCollector.Collect(c.Message())

	var (
		costPerImpression float64
		done              bool
		resp              intele.Response
	)
	for {
		resp, err = h.inputManager.Get(context.Background(), c.Sender().ID, 0)
		if resp.Message != nil {
			inputCollector.Collect(resp.Message)
		}
		switch {
		case resp.Canceled:
			_ = inputCollector.Clear(c, collector.ClearOptions{IgnoreErrors: true, ExcludeLast: true})
			return nil
		case err != nil:
			log.Error().Msgf("(user: %d) error while input cost per impression: %v", c.Sender().ID, err)
			_ = inputCollector.Send(c,
				h.banners.Advertiser.Caption(h.layout.Text(c, "input_error", h.layout.Text(c, "input_cost_per_impression"))),
				h.layout.Markup(c, "campaign:update:back", struct {
					ID   string
					Page int
				}{
					ID:   campaignID,
					Page: p,
				}),
			)
		case resp.Message == nil:
			log.Error().Msgf("(user: %d) error while input cost per impression: %v", c.Sender().ID, err)
			_ = inputCollector.Send(c,
				h.banners.Advertiser.Caption(h.layout.Text(c, "input_error", h.layout.Text(c, "input_cost_per_impression"))),
				h.layout.Markup(c, "campaign:update:back", struct {
					ID   string
					Page int
				}{
					ID:   campaignID,
					Page: p,
				}),
			)
		case !validator.CostPerImpression(resp.Message.Text, nil):
			_ = inputCollector.Send(c,
				h.banners.Advertiser.Caption(h.layout.Text(c, "invalid_cost_per_impression")),
				h.layout.Markup(c, "campaign:update:back", struct {
					ID   string
					Page int
				}{
					ID:   campaignID,
					Page: p,
				}),
			)
		case validator.CostPerImpression(resp.Message.Text, nil):
			costPerImpression, _ = strconv.ParseFloat(resp.Message.Text, 64)
			_ = inputCollector.Clear(c, collector.ClearOptions{IgnoreErrors: true})
			done = true
		}
		if done {
			break
		}
	}

	updateRequestBody := request.UpdateCampaign{
		ImpressionsLimit:  campaign.ImpressionsLimit,
		ClicksLimit:       campaign.ClicksLimit,
		CostPerImpression: costPerImpression,
		CostPerClick:      campaign.CostPerClick,
		AdTitle:           &campaign.AdTitle,
		AdText:            &campaign.AdText,
		StartDate:         campaign.StartDate,
		EndDate:           campaign.EndDate,
		Targeting: &request.Targeting{
			Gender:   campaign.Targeting.Gender,
			AgeFrom:  campaign.Targeting.AgeFrom,
			AgeTo:    campaign.Targeting.AgeTo,
			Location: campaign.Targeting.Location,
		},
	}

	updatedCampaign, err := h.advertisingClient.UpdateCampaign(updateRequestBody, user.PlatformID, campaignID)
	switch {
	case errors.Is(err, errorz.ErrInvalidData) || errors.Is(err, errorz.ErrCampaignNotFound):
		return c.Send(
			h.banners.Advertiser.Caption(h.layout.Text(c, "invalid_cost_per_impression")),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	case err != nil:
		return c.Send(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	return c.Send(
		h.banners.Advertiser.Caption(h.layout.Text(c, "campaign_updated_successfully", struct {
			AdTitle string
		}{
			AdTitle: updatedCampaign.AdTitle,
		})),
		h.layout.Markup(c, "campaign:update:back", struct {
			ID   string
			Page int
		}{
			ID:   campaignID,
			Page: p,
		}),
	)
}

func (h *Handler) updateCostPerClick(c tele.Context) error {
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
			h.layout.Markup(c, "campaign:update:back", struct {
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
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	currentDate, err := h.advertisingClient.GetCurrentDate()
	if err != nil {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	if campaign.IsOver(currentDate) {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "campaign_is_over")),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	inputCollector := collector.New()
	_ = c.Edit(
		h.banners.Advertiser.Caption(h.layout.Text(c, "input_cost_per_click")),
		h.layout.Markup(c, "campaign:update:back", struct {
			ID   string
			Page int
		}{
			ID:   campaignID,
			Page: p,
		}),
	)
	inputCollector.Collect(c.Message())

	var (
		costPerClick float64
		done         bool
		resp         intele.Response
	)
	for {
		resp, err = h.inputManager.Get(context.Background(), c.Sender().ID, 0)
		if resp.Message != nil {
			inputCollector.Collect(resp.Message)
		}
		switch {
		case resp.Canceled:
			_ = inputCollector.Clear(c, collector.ClearOptions{IgnoreErrors: true, ExcludeLast: true})
			return nil
		case err != nil:
			log.Error().Msgf("(user: %d) error while input cost per click: %v", c.Sender().ID, err)
			_ = inputCollector.Send(c,
				h.banners.Advertiser.Caption(h.layout.Text(c, "input_error", h.layout.Text(c, "input_cost_per_click"))),
				h.layout.Markup(c, "campaign:update:back", struct {
					ID   string
					Page int
				}{
					ID:   campaignID,
					Page: p,
				}),
			)
		case resp.Message == nil:
			log.Error().Msgf("(user: %d) error while input cost per click: %v", c.Sender().ID, err)
			_ = inputCollector.Send(c,
				h.banners.Advertiser.Caption(h.layout.Text(c, "input_error", h.layout.Text(c, "input_cost_per_click"))),
				h.layout.Markup(c, "campaign:update:back", struct {
					ID   string
					Page int
				}{
					ID:   campaignID,
					Page: p,
				}),
			)
		case !validator.CostPerClick(resp.Message.Text, nil):
			_ = inputCollector.Send(c,
				h.banners.Advertiser.Caption(h.layout.Text(c, "invalid_cost_per_click")),
				h.layout.Markup(c, "campaign:update:back", struct {
					ID   string
					Page int
				}{
					ID:   campaignID,
					Page: p,
				}),
			)
		case validator.CostPerClick(resp.Message.Text, nil):
			costPerClick, _ = strconv.ParseFloat(resp.Message.Text, 64)
			_ = inputCollector.Clear(c, collector.ClearOptions{IgnoreErrors: true})
			done = true
		}
		if done {
			break
		}
	}

	updateRequestBody := request.UpdateCampaign{
		ImpressionsLimit:  campaign.ImpressionsLimit,
		ClicksLimit:       campaign.ClicksLimit,
		CostPerImpression: campaign.CostPerImpression,
		CostPerClick:      costPerClick,
		AdTitle:           &campaign.AdTitle,
		AdText:            &campaign.AdText,
		StartDate:         campaign.StartDate,
		EndDate:           campaign.EndDate,
		Targeting: &request.Targeting{
			Gender:   campaign.Targeting.Gender,
			AgeFrom:  campaign.Targeting.AgeFrom,
			AgeTo:    campaign.Targeting.AgeTo,
			Location: campaign.Targeting.Location,
		},
	}

	updatedCampaign, err := h.advertisingClient.UpdateCampaign(updateRequestBody, user.PlatformID, campaignID)
	switch {
	case errors.Is(err, errorz.ErrInvalidData) || errors.Is(err, errorz.ErrCampaignNotFound):
		return c.Send(
			h.banners.Advertiser.Caption(h.layout.Text(c, "invalid_cost_per_click")),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	case err != nil:
		return c.Send(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	return c.Send(
		h.banners.Advertiser.Caption(h.layout.Text(c, "campaign_updated_successfully", struct {
			AdTitle string
		}{
			AdTitle: updatedCampaign.AdTitle,
		})),
		h.layout.Markup(c, "campaign:update:back", struct {
			ID   string
			Page int
		}{
			ID:   campaignID,
			Page: p,
		}),
	)
}

func (h *Handler) updateAdTitle(c tele.Context) error {
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
			h.layout.Markup(c, "campaign:update:back", struct {
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
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	currentDate, err := h.advertisingClient.GetCurrentDate()
	if err != nil {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	if campaign.IsOver(currentDate) {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "campaign_is_over")),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	markup := h.layout.Markup(c, "campaign:update:back", struct {
		ID   string
		Page int
	}{
		ID:   campaignID,
		Page: p,
	})
	markup.InlineKeyboard = append(
		[][]tele.InlineButton{{*h.layout.Button(c, "update:ad_title:set_empty").Inline()}},
		markup.InlineKeyboard...,
	)

	inputCollector := collector.New()
	_ = c.Edit(
		h.banners.Advertiser.Caption(h.layout.Text(c, "input_ad_title")),
		markup,
	)
	inputCollector.Collect(c.Message())

	var (
		adTitle *string
		done    bool
		resp    intele.Response
	)
	for {
		resp, err = h.inputManager.Get(context.Background(), c.Sender().ID, 0, h.layout.Button(c, "update:ad_title:set_empty"))
		if resp.Message != nil {
			inputCollector.Collect(resp.Message)
		}
		switch {
		case resp.Canceled:
			_ = inputCollector.Clear(c, collector.ClearOptions{IgnoreErrors: true, ExcludeLast: true})
			return nil
		case err != nil:
			log.Error().Msgf("(user: %d) error while input ad title: %v", c.Sender().ID, err)
			_ = inputCollector.Send(c,
				h.banners.Advertiser.Caption(h.layout.Text(c, "input_error", h.layout.Text(c, "input_ad_title"))),
				h.layout.Markup(c, "campaign:update:back", struct {
					ID   string
					Page int
				}{
					ID:   campaignID,
					Page: p,
				}),
			)
		case resp.Callback != nil:
			adTitle = nil
			_ = inputCollector.Clear(c, collector.ClearOptions{IgnoreErrors: true})
			done = true
		case !validator.AdTitle(resp.Message.Text, nil):
			_ = inputCollector.Send(c,
				h.banners.Advertiser.Caption(h.layout.Text(c, "invalid_ad_title")),
				h.layout.Markup(c, "campaign:update:back", struct {
					ID   string
					Page int
				}{
					ID:   campaignID,
					Page: p,
				}),
			)
		case validator.AdTitle(resp.Message.Text, nil):
			adTitle = &resp.Message.Text
			_ = inputCollector.Clear(c, collector.ClearOptions{IgnoreErrors: true})
			done = true
		}
		if done {
			break
		}
	}

	updateRequestBody := request.UpdateCampaign{
		ImpressionsLimit:  campaign.ImpressionsLimit,
		ClicksLimit:       campaign.ClicksLimit,
		CostPerImpression: campaign.CostPerImpression,
		CostPerClick:      campaign.CostPerClick,
		AdTitle:           adTitle,
		AdText:            &campaign.AdText,
		StartDate:         campaign.StartDate,
		EndDate:           campaign.EndDate,
		Targeting: &request.Targeting{
			Gender:   campaign.Targeting.Gender,
			AgeFrom:  campaign.Targeting.AgeFrom,
			AgeTo:    campaign.Targeting.AgeTo,
			Location: campaign.Targeting.Location,
		},
	}

	updatedCampaign, err := h.advertisingClient.UpdateCampaign(updateRequestBody, user.PlatformID, campaignID)
	switch {
	case errors.Is(err, errorz.ErrInvalidData) || errors.Is(err, errorz.ErrCampaignNotFound):
		return c.Send(
			h.banners.Advertiser.Caption(h.layout.Text(c, "invalid_ad_title")),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	case err != nil:
		return c.Send(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	return c.Send(
		h.banners.Advertiser.Caption(h.layout.Text(c, "campaign_updated_successfully", struct {
			AdTitle string
		}{
			AdTitle: updatedCampaign.AdTitle,
		})),
		h.layout.Markup(c, "campaign:update:back", struct {
			ID   string
			Page int
		}{
			ID:   campaignID,
			Page: p,
		}),
	)
}

func (h *Handler) updateAdText(c tele.Context) error {
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
			h.layout.Markup(c, "campaign:update:back", struct {
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
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	currentDate, err := h.advertisingClient.GetCurrentDate()
	if err != nil {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	if campaign.IsOver(currentDate) {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "campaign_is_over")),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	markup := h.layout.Markup(c, "campaign:update:back", struct {
		ID   string
		Page int
	}{
		ID:   campaignID,
		Page: p,
	})
	markup.InlineKeyboard = append(
		[][]tele.InlineButton{{*h.layout.Button(c, "update:ad_text:set_empty").Inline()}},
		markup.InlineKeyboard...,
	)

	inputCollector := collector.New()
	_ = c.Edit(
		h.banners.Advertiser.Caption(h.layout.Text(c, "input_ad_text")),
		markup,
	)
	inputCollector.Collect(c.Message())

	var (
		adText *string
		done   bool
		resp   intele.Response
	)
	for {
		resp, err = h.inputManager.Get(context.Background(), c.Sender().ID, 0, h.layout.Button(c, "update:ad_text:set_empty"))
		if resp.Message != nil {
			inputCollector.Collect(resp.Message)
		}
		switch {
		case resp.Canceled:
			_ = inputCollector.Clear(c, collector.ClearOptions{IgnoreErrors: true, ExcludeLast: true})
			return nil
		case err != nil:
			log.Error().Msgf("(user: %d) error while input ad text: %v", c.Sender().ID, err)
			_ = inputCollector.Send(c,
				h.banners.Advertiser.Caption(h.layout.Text(c, "input_error", h.layout.Text(c, "input_ad_text"))),
				h.layout.Markup(c, "campaign:update:back", struct {
					ID   string
					Page int
				}{
					ID:   campaignID,
					Page: p,
				}),
			)
		case resp.Callback != nil:
			adText = nil
			_ = inputCollector.Clear(c, collector.ClearOptions{IgnoreErrors: true})
			done = true
		case !validator.AdText(resp.Message.Text, nil):
			_ = inputCollector.Send(c,
				h.banners.Advertiser.Caption(h.layout.Text(c, "invalid_ad_text")),
				h.layout.Markup(c, "campaign:update:back", struct {
					ID   string
					Page int
				}{
					ID:   campaignID,
					Page: p,
				}),
			)
		case validator.AdText(resp.Message.Text, nil):
			adText = &resp.Message.Text
			_ = inputCollector.Clear(c, collector.ClearOptions{IgnoreErrors: true})
			done = true
		}
		if done {
			break
		}
	}

	updateRequestBody := request.UpdateCampaign{
		ImpressionsLimit:  campaign.ImpressionsLimit,
		ClicksLimit:       campaign.ClicksLimit,
		CostPerImpression: campaign.CostPerImpression,
		CostPerClick:      campaign.CostPerClick,
		AdTitle:           &campaign.AdTitle,
		AdText:            adText,
		StartDate:         campaign.StartDate,
		EndDate:           campaign.EndDate,
		Targeting: &request.Targeting{
			Gender:   campaign.Targeting.Gender,
			AgeFrom:  campaign.Targeting.AgeFrom,
			AgeTo:    campaign.Targeting.AgeTo,
			Location: campaign.Targeting.Location,
		},
	}

	updatedCampaign, err := h.advertisingClient.UpdateCampaign(updateRequestBody, user.PlatformID, campaignID)
	switch {
	case errors.Is(err, errorz.ErrInvalidData) || errors.Is(err, errorz.ErrCampaignNotFound):
		return c.Send(
			h.banners.Advertiser.Caption(h.layout.Text(c, "invalid_ad_text")),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	case err != nil:
		return c.Send(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	return c.Send(
		h.banners.Advertiser.Caption(h.layout.Text(c, "campaign_updated_successfully", struct {
			AdTitle string
		}{
			AdTitle: updatedCampaign.AdTitle,
		})),
		h.layout.Markup(c, "campaign:update:back", struct {
			ID   string
			Page int
		}{
			ID:   campaignID,
			Page: p,
		}),
	)
}

func (h *Handler) updateStartDate(c tele.Context) error {
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
			h.layout.Markup(c, "campaign:update:back", struct {
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
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	currentDate, err := h.advertisingClient.GetCurrentDate()
	if err != nil {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	if campaign.IsActive(currentDate) {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "campaign_is_active")),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}
	if campaign.IsOver(currentDate) {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "campaign_is_over")),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	params := make(map[string]interface{})
	params["currentDate"] = currentDate

	inputCollector := collector.New()
	_ = c.Edit(
		h.banners.Advertiser.Caption(h.layout.Text(c, "input_start_date", struct {
			CurrentDate int
		}{
			CurrentDate: currentDate,
		})),
		h.layout.Markup(c, "campaign:update:back", struct {
			ID   string
			Page int
		}{
			ID:   campaignID,
			Page: p,
		}),
	)
	inputCollector.Collect(c.Message())

	var (
		startDate int
		done      bool
		resp      intele.Response
	)
	for {
		resp, err = h.inputManager.Get(context.Background(), c.Sender().ID, 0)
		if resp.Message != nil {
			inputCollector.Collect(resp.Message)
		}
		switch {
		case resp.Canceled:
			_ = inputCollector.Clear(c, collector.ClearOptions{IgnoreErrors: true, ExcludeLast: true})
			return nil
		case err != nil:
			log.Error().Msgf("(user: %d) error while input start date: %v", c.Sender().ID, err)
			_ = inputCollector.Send(c,
				h.banners.Advertiser.Caption(h.layout.Text(c, "input_error", h.layout.Text(c, "input_start_date", struct {
					CurrentDate int
				}{
					CurrentDate: currentDate,
				}))),
				h.layout.Markup(c, "campaign:update:back", struct {
					ID   string
					Page int
				}{
					ID:   campaignID,
					Page: p,
				}),
			)
		case resp.Message == nil:
			log.Error().Msgf("(user: %d) error while input start date: %v", c.Sender().ID, err)
			_ = inputCollector.Send(c,
				h.banners.Advertiser.Caption(h.layout.Text(c, "input_error", h.layout.Text(c, "input_start_date", struct {
					CurrentDate int
				}{
					CurrentDate: currentDate,
				}))),
				h.layout.Markup(c, "campaign:update:back", struct {
					ID   string
					Page int
				}{
					ID:   campaignID,
					Page: p,
				}),
			)
		case !validator.StartDate(resp.Message.Text, params):
			_ = inputCollector.Send(c,
				h.banners.Advertiser.Caption(h.layout.Text(c, "invalid_start_date", struct {
					CurrentDate int
				}{
					CurrentDate: currentDate,
				})),
				h.layout.Markup(c, "campaign:update:back", struct {
					ID   string
					Page int
				}{
					ID:   campaignID,
					Page: p,
				}),
			)
		case validator.StartDate(resp.Message.Text, params):
			startDate, _ = strconv.Atoi(resp.Message.Text)
			_ = inputCollector.Clear(c, collector.ClearOptions{IgnoreErrors: true})
			done = true
		}
		if done {
			break
		}
	}

	updateRequestBody := request.UpdateCampaign{
		ImpressionsLimit:  campaign.ImpressionsLimit,
		ClicksLimit:       campaign.ClicksLimit,
		CostPerImpression: campaign.CostPerImpression,
		CostPerClick:      campaign.CostPerClick,
		AdTitle:           &campaign.AdTitle,
		AdText:            &campaign.AdText,
		StartDate:         startDate,
		EndDate:           campaign.EndDate,
		Targeting: &request.Targeting{
			Gender:   campaign.Targeting.Gender,
			AgeFrom:  campaign.Targeting.AgeFrom,
			AgeTo:    campaign.Targeting.AgeTo,
			Location: campaign.Targeting.Location,
		},
	}

	updatedCampaign, err := h.advertisingClient.UpdateCampaign(updateRequestBody, user.PlatformID, campaignID)
	switch {
	case errors.Is(err, errorz.ErrInvalidData) || errors.Is(err, errorz.ErrCampaignNotFound):
		return c.Send(
			h.banners.Advertiser.Caption(h.layout.Text(c, "invalid_start_date", struct {
				CurrentDate int
			}{
				CurrentDate: currentDate,
			})),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	case err != nil:
		return c.Send(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	return c.Send(
		h.banners.Advertiser.Caption(h.layout.Text(c, "campaign_updated_successfully", struct {
			AdTitle string
		}{
			AdTitle: updatedCampaign.AdTitle,
		})),
		h.layout.Markup(c, "campaign:update:back", struct {
			ID   string
			Page int
		}{
			ID:   campaignID,
			Page: p,
		}),
	)
}

func (h *Handler) updateEndDate(c tele.Context) error {
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
			h.layout.Markup(c, "campaign:update:back", struct {
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
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	currentDate, err := h.advertisingClient.GetCurrentDate()
	if err != nil {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	if campaign.IsActive(currentDate) {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "campaign_is_active")),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}
	if campaign.IsOver(currentDate) {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "campaign_is_over")),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	params := make(map[string]interface{})
	params["startDate"] = campaign.StartDate

	inputCollector := collector.New()
	_ = c.Edit(
		h.banners.Advertiser.Caption(h.layout.Text(c, "input_end_date", struct {
			StartDate int
		}{
			StartDate: campaign.StartDate,
		})),
		h.layout.Markup(c, "campaign:update:back", struct {
			ID   string
			Page int
		}{
			ID:   campaignID,
			Page: p,
		}),
	)
	inputCollector.Collect(c.Message())

	var (
		endDate int
		done    bool
		resp    intele.Response
	)
	for {
		resp, err = h.inputManager.Get(context.Background(), c.Sender().ID, 0)
		if resp.Message != nil {
			inputCollector.Collect(resp.Message)
		}
		switch {
		case resp.Canceled:
			_ = inputCollector.Clear(c, collector.ClearOptions{IgnoreErrors: true, ExcludeLast: true})
			return nil
		case err != nil:
			log.Error().Msgf("(user: %d) error while input end date: %v", c.Sender().ID, err)
			_ = inputCollector.Send(c,
				h.banners.Advertiser.Caption(h.layout.Text(c, "input_error", h.layout.Text(c, "input_end_date", struct {
					StartDate int
				}{
					StartDate: campaign.StartDate,
				}))),
				h.layout.Markup(c, "campaign:update:back", struct {
					ID   string
					Page int
				}{
					ID:   campaignID,
					Page: p,
				}),
			)
		case resp.Message == nil:
			log.Error().Msgf("(user: %d) error while input end date: %v", c.Sender().ID, err)
			_ = inputCollector.Send(c,
				h.banners.Advertiser.Caption(h.layout.Text(c, "input_error", h.layout.Text(c, "input_end_date", struct {
					StartDate int
				}{
					StartDate: campaign.StartDate,
				}))),
				h.layout.Markup(c, "campaign:update:back", struct {
					ID   string
					Page int
				}{
					ID:   campaignID,
					Page: p,
				}),
			)
		case !validator.EndDate(resp.Message.Text, params):
			_ = inputCollector.Send(c,
				h.banners.Advertiser.Caption(h.layout.Text(c, "invalid_end_date", struct {
					StartDate int
				}{
					StartDate: campaign.StartDate,
				})),
				h.layout.Markup(c, "campaign:update:back", struct {
					ID   string
					Page int
				}{
					ID:   campaignID,
					Page: p,
				}),
			)
		case validator.EndDate(resp.Message.Text, params):
			endDate, _ = strconv.Atoi(resp.Message.Text)
			_ = inputCollector.Clear(c, collector.ClearOptions{IgnoreErrors: true})
			done = true
		}
		if done {
			break
		}
	}

	updateRequestBody := request.UpdateCampaign{
		ImpressionsLimit:  campaign.ImpressionsLimit,
		ClicksLimit:       campaign.ClicksLimit,
		CostPerImpression: campaign.CostPerImpression,
		CostPerClick:      campaign.CostPerClick,
		AdTitle:           &campaign.AdTitle,
		AdText:            &campaign.AdText,
		StartDate:         campaign.StartDate,
		EndDate:           endDate,
		Targeting: &request.Targeting{
			Gender:   campaign.Targeting.Gender,
			AgeFrom:  campaign.Targeting.AgeFrom,
			AgeTo:    campaign.Targeting.AgeTo,
			Location: campaign.Targeting.Location,
		},
	}

	updatedCampaign, err := h.advertisingClient.UpdateCampaign(updateRequestBody, user.PlatformID, campaignID)
	switch {
	case errors.Is(err, errorz.ErrInvalidData) || errors.Is(err, errorz.ErrCampaignNotFound):
		return c.Send(
			h.banners.Advertiser.Caption(h.layout.Text(c, "invalid_end_date", struct {
				StartDate int
			}{
				StartDate: campaign.StartDate,
			})),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	case err != nil:
		return c.Send(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	return c.Send(
		h.banners.Advertiser.Caption(h.layout.Text(c, "campaign_updated_successfully", struct {
			AdTitle string
		}{
			AdTitle: updatedCampaign.AdTitle,
		})),
		h.layout.Markup(c, "campaign:update:back", struct {
			ID   string
			Page int
		}{
			ID:   campaignID,
			Page: p,
		}),
	)
}

func (h *Handler) updateGender(c tele.Context) error {
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
			h.layout.Markup(c, "campaign:update:back", struct {
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
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	currentDate, err := h.advertisingClient.GetCurrentDate()
	if err != nil {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	if campaign.IsOver(currentDate) {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "campaign_is_over")),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	markup := h.layout.Markup(c, "campaign:update:back", struct {
		ID   string
		Page int
	}{
		ID:   campaignID,
		Page: p,
	})
	markup.InlineKeyboard = append(
		[][]tele.InlineButton{{*h.layout.Button(c, "update:gender:set_empty").Inline()}},
		markup.InlineKeyboard...,
	)

	inputCollector := collector.New()
	_ = c.Edit(
		h.banners.Advertiser.Caption(h.layout.Text(c, "input_targeting_gender")),
		markup,
	)
	inputCollector.Collect(c.Message())

	var (
		gender *string
		done   bool
		resp   intele.Response
	)
	for {
		resp, err = h.inputManager.Get(context.Background(), c.Sender().ID, 0, h.layout.Button(c, "update:gender:set_empty"))
		if resp.Message != nil {
			inputCollector.Collect(resp.Message)
		}
		switch {
		case resp.Canceled:
			_ = inputCollector.Clear(c, collector.ClearOptions{IgnoreErrors: true, ExcludeLast: true})
			return nil
		case err != nil:
			log.Error().Msgf("(user: %d) error while input targeting gender: %v", c.Sender().ID, err)
			_ = inputCollector.Send(c,
				h.banners.Advertiser.Caption(h.layout.Text(c, "input_error", h.layout.Text(c, "input_targeting_gender"))),
				h.layout.Markup(c, "campaign:update:back", struct {
					ID   string
					Page int
				}{
					ID:   campaignID,
					Page: p,
				}),
			)
		case resp.Callback != nil:
			gender = nil
			_ = inputCollector.Clear(c, collector.ClearOptions{IgnoreErrors: true})
			done = true
		case !validator.TargetingGender(resp.Message.Text, nil):
			_ = inputCollector.Send(c,
				h.banners.Advertiser.Caption(h.layout.Text(c, "invalid_targeting_gender")),
				h.layout.Markup(c, "campaign:update:back", struct {
					ID   string
					Page int
				}{
					ID:   campaignID,
					Page: p,
				}),
			)
		case validator.TargetingGender(resp.Message.Text, nil):
			gender = &resp.Message.Text
			_ = inputCollector.Clear(c, collector.ClearOptions{IgnoreErrors: true})
			done = true
		}
		if done {
			break
		}
	}

	updateRequestBody := request.UpdateCampaign{
		ImpressionsLimit:  campaign.ImpressionsLimit,
		ClicksLimit:       campaign.ClicksLimit,
		CostPerImpression: campaign.CostPerImpression,
		CostPerClick:      campaign.CostPerClick,
		AdTitle:           &campaign.AdTitle,
		AdText:            &campaign.AdText,
		StartDate:         campaign.StartDate,
		EndDate:           campaign.EndDate,
		Targeting: &request.Targeting{
			Gender:   gender,
			AgeFrom:  campaign.Targeting.AgeFrom,
			AgeTo:    campaign.Targeting.AgeTo,
			Location: campaign.Targeting.Location,
		},
	}

	updatedCampaign, err := h.advertisingClient.UpdateCampaign(updateRequestBody, user.PlatformID, campaignID)
	switch {
	case errors.Is(err, errorz.ErrInvalidData) || errors.Is(err, errorz.ErrCampaignNotFound):
		return c.Send(
			h.banners.Advertiser.Caption(h.layout.Text(c, "invalid_targeting_gender")),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	case err != nil:
		return c.Send(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	return c.Send(
		h.banners.Advertiser.Caption(h.layout.Text(c, "campaign_updated_successfully", struct {
			AdTitle string
		}{
			AdTitle: updatedCampaign.AdTitle,
		})),
		h.layout.Markup(c, "campaign:update:back", struct {
			ID   string
			Page int
		}{
			ID:   campaignID,
			Page: p,
		}),
	)
}

func (h *Handler) updateLocation(c tele.Context) error {
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
			h.layout.Markup(c, "campaign:update:back", struct {
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
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	currentDate, err := h.advertisingClient.GetCurrentDate()
	if err != nil {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	if campaign.IsOver(currentDate) {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "campaign_is_over")),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	markup := h.layout.Markup(c, "campaign:update:back", struct {
		ID   string
		Page int
	}{
		ID:   campaignID,
		Page: p,
	})
	markup.InlineKeyboard = append(
		[][]tele.InlineButton{{*h.layout.Button(c, "update:location:set_empty").Inline()}},
		markup.InlineKeyboard...,
	)

	inputCollector := collector.New()
	_ = c.Edit(
		h.banners.Advertiser.Caption(h.layout.Text(c, "input_targeting_location")),
		markup,
	)
	inputCollector.Collect(c.Message())

	var (
		location *string
		done     bool
		resp     intele.Response
	)
	for {
		resp, err = h.inputManager.Get(context.Background(), c.Sender().ID, 0, h.layout.Button(c, "update:location:set_empty"))
		if resp.Message != nil {
			inputCollector.Collect(resp.Message)
		}
		switch {
		case resp.Canceled:
			_ = inputCollector.Clear(c, collector.ClearOptions{IgnoreErrors: true, ExcludeLast: true})
			return nil
		case err != nil:
			log.Error().Msgf("(user: %d) error while input targeting location: %v", c.Sender().ID, err)
			_ = inputCollector.Send(c,
				h.banners.Advertiser.Caption(h.layout.Text(c, "input_error", h.layout.Text(c, "input_targeting_location"))),
				h.layout.Markup(c, "campaign:update:back", struct {
					ID   string
					Page int
				}{
					ID:   campaignID,
					Page: p,
				}),
			)
		case resp.Callback != nil:
			location = nil
			_ = inputCollector.Clear(c, collector.ClearOptions{IgnoreErrors: true})
			done = true
		case !validator.TargetingLocation(resp.Message.Text, nil):
			_ = inputCollector.Send(c,
				h.banners.Advertiser.Caption(h.layout.Text(c, "invalid_targeting_location")),
				h.layout.Markup(c, "campaign:update:back", struct {
					ID   string
					Page int
				}{
					ID:   campaignID,
					Page: p,
				}),
			)
		case validator.TargetingLocation(resp.Message.Text, nil):
			location = &resp.Message.Text
			_ = inputCollector.Clear(c, collector.ClearOptions{IgnoreErrors: true})
			done = true
		}
		if done {
			break
		}
	}

	updateRequestBody := request.UpdateCampaign{
		ImpressionsLimit:  campaign.ImpressionsLimit,
		ClicksLimit:       campaign.ClicksLimit,
		CostPerImpression: campaign.CostPerImpression,
		CostPerClick:      campaign.CostPerClick,
		AdTitle:           &campaign.AdTitle,
		AdText:            &campaign.AdText,
		StartDate:         campaign.StartDate,
		EndDate:           campaign.EndDate,
		Targeting: &request.Targeting{
			Gender:   campaign.Targeting.Gender,
			AgeFrom:  campaign.Targeting.AgeFrom,
			AgeTo:    campaign.Targeting.AgeTo,
			Location: location,
		},
	}

	updatedCampaign, err := h.advertisingClient.UpdateCampaign(updateRequestBody, user.PlatformID, campaignID)
	switch {
	case errors.Is(err, errorz.ErrInvalidData) || errors.Is(err, errorz.ErrCampaignNotFound):
		return c.Send(
			h.banners.Advertiser.Caption(h.layout.Text(c, "invalid_targeting_location")),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	case err != nil:
		return c.Send(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	return c.Send(
		h.banners.Advertiser.Caption(h.layout.Text(c, "campaign_updated_successfully", struct {
			AdTitle string
		}{
			AdTitle: updatedCampaign.AdTitle,
		})),
		h.layout.Markup(c, "campaign:update:back", struct {
			ID   string
			Page int
		}{
			ID:   campaignID,
			Page: p,
		}),
	)
}

func (h *Handler) updateAgeFrom(c tele.Context) error {
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
			h.layout.Markup(c, "campaign:update:back", struct {
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
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	currentDate, err := h.advertisingClient.GetCurrentDate()
	if err != nil {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	if campaign.IsOver(currentDate) {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "campaign_is_over")),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	markup := h.layout.Markup(c, "campaign:update:back", struct {
		ID   string
		Page int
	}{
		ID:   campaignID,
		Page: p,
	})
	markup.InlineKeyboard = append(
		[][]tele.InlineButton{{*h.layout.Button(c, "update:age_from:set_empty").Inline()}},
		markup.InlineKeyboard...,
	)

	inputCollector := collector.New()
	_ = c.Edit(
		h.banners.Advertiser.Caption(h.layout.Text(c, "input_targeting_age_from")),
		markup,
	)
	inputCollector.Collect(c.Message())

	var (
		ageFrom *int
		done    bool
		resp    intele.Response
	)
	for {
		resp, err = h.inputManager.Get(context.Background(), c.Sender().ID, 0, h.layout.Button(c, "update:age_from:set_empty"))
		if resp.Message != nil {
			inputCollector.Collect(resp.Message)
		}
		switch {
		case resp.Canceled:
			_ = inputCollector.Clear(c, collector.ClearOptions{IgnoreErrors: true, ExcludeLast: true})
			return nil
		case err != nil:
			log.Error().Msgf("(user: %d) error while input targeting age from: %v", c.Sender().ID, err)
			_ = inputCollector.Send(c,
				h.banners.Advertiser.Caption(h.layout.Text(c, "input_error", h.layout.Text(c, "input_targeting_age_from"))),
				h.layout.Markup(c, "campaign:update:back", struct {
					ID   string
					Page int
				}{
					ID:   campaignID,
					Page: p,
				}),
			)
		case resp.Callback != nil:
			ageFrom = nil
			_ = inputCollector.Clear(c, collector.ClearOptions{IgnoreErrors: true})
			done = true
		case !validator.TargetingAgeFrom(resp.Message.Text, nil):
			_ = inputCollector.Send(c,
				h.banners.Advertiser.Caption(h.layout.Text(c, "invalid_targeting_age_from")),
				h.layout.Markup(c, "campaign:update:back", struct {
					ID   string
					Page int
				}{
					ID:   campaignID,
					Page: p,
				}),
			)
		case validator.TargetingAgeFrom(resp.Message.Text, nil):
			temp, _ := strconv.Atoi(resp.Message.Text)
			ageFrom = &temp
			_ = inputCollector.Clear(c, collector.ClearOptions{IgnoreErrors: true})
			done = true
		}
		if done {
			break
		}
	}

	updateRequestBody := request.UpdateCampaign{
		ImpressionsLimit:  campaign.ImpressionsLimit,
		ClicksLimit:       campaign.ClicksLimit,
		CostPerImpression: campaign.CostPerImpression,
		CostPerClick:      campaign.CostPerClick,
		AdTitle:           &campaign.AdTitle,
		AdText:            &campaign.AdText,
		StartDate:         campaign.StartDate,
		EndDate:           campaign.EndDate,
		Targeting: &request.Targeting{
			Gender:   campaign.Targeting.Gender,
			AgeFrom:  ageFrom,
			AgeTo:    campaign.Targeting.AgeTo,
			Location: campaign.Targeting.Location,
		},
	}

	updatedCampaign, err := h.advertisingClient.UpdateCampaign(updateRequestBody, user.PlatformID, campaignID)
	switch {
	case errors.Is(err, errorz.ErrInvalidData) || errors.Is(err, errorz.ErrCampaignNotFound):
		return c.Send(
			h.banners.Advertiser.Caption(h.layout.Text(c, "invalid_targeting_age_from")),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	case err != nil:
		return c.Send(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	return c.Send(
		h.banners.Advertiser.Caption(h.layout.Text(c, "campaign_updated_successfully", struct {
			AdTitle string
		}{
			AdTitle: updatedCampaign.AdTitle,
		})),
		h.layout.Markup(c, "campaign:update:back", struct {
			ID   string
			Page int
		}{
			ID:   campaignID,
			Page: p,
		}),
	)
}

func (h *Handler) updateAgeTo(c tele.Context) error {
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
			h.layout.Markup(c, "campaign:update:back", struct {
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
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	currentDate, err := h.advertisingClient.GetCurrentDate()
	if err != nil {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	if campaign.IsOver(currentDate) {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "campaign_is_over")),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	markup := h.layout.Markup(c, "campaign:update:back", struct {
		ID   string
		Page int
	}{
		ID:   campaignID,
		Page: p,
	})
	markup.InlineKeyboard = append(
		[][]tele.InlineButton{{*h.layout.Button(c, "update:age_to:set_empty").Inline()}},
		markup.InlineKeyboard...,
	)

	params := make(map[string]interface{})
	params["ageFrom"] = 0
	if campaign.Targeting.AgeFrom != nil {
		params["ageFrom"] = *campaign.Targeting.AgeFrom
	}

	inputCollector := collector.New()
	_ = c.Edit(
		h.banners.Advertiser.Caption(h.layout.Text(c, "input_targeting_age_to")),
		markup,
	)
	inputCollector.Collect(c.Message())

	var (
		ageTo *int
		done  bool
		resp  intele.Response
	)
	for {
		resp, err = h.inputManager.Get(context.Background(), c.Sender().ID, 0, h.layout.Button(c, "update:age_to:set_empty"))
		if resp.Message != nil {
			inputCollector.Collect(resp.Message)
		}
		switch {
		case resp.Canceled:
			_ = inputCollector.Clear(c, collector.ClearOptions{IgnoreErrors: true, ExcludeLast: true})
			return nil
		case err != nil:
			log.Error().Msgf("(user: %d) error while input targeting age to: %v", c.Sender().ID, err)
			_ = inputCollector.Send(c,
				h.banners.Advertiser.Caption(h.layout.Text(c, "input_error", h.layout.Text(c, "input_targeting_age_to"))),
				h.layout.Markup(c, "campaign:update:back", struct {
					ID   string
					Page int
				}{
					ID:   campaignID,
					Page: p,
				}),
			)
		case resp.Callback != nil:
			ageTo = nil
			_ = inputCollector.Clear(c, collector.ClearOptions{IgnoreErrors: true})
			done = true
		case !validator.TargetingAgeTo(resp.Message.Text, params):
			_ = inputCollector.Send(c,
				h.banners.Advertiser.Caption(h.layout.Text(c, "invalid_targeting_age_to")),
				h.layout.Markup(c, "campaign:update:back", struct {
					ID   string
					Page int
				}{
					ID:   campaignID,
					Page: p,
				}),
			)
		case validator.TargetingAgeTo(resp.Message.Text, params):
			temp, _ := strconv.Atoi(resp.Message.Text)
			ageTo = &temp
			_ = inputCollector.Clear(c, collector.ClearOptions{IgnoreErrors: true})
			done = true
		}
		if done {
			break
		}
	}

	updateRequestBody := request.UpdateCampaign{
		ImpressionsLimit:  campaign.ImpressionsLimit,
		ClicksLimit:       campaign.ClicksLimit,
		CostPerImpression: campaign.CostPerImpression,
		CostPerClick:      campaign.CostPerClick,
		AdTitle:           &campaign.AdTitle,
		AdText:            &campaign.AdText,
		StartDate:         campaign.StartDate,
		EndDate:           campaign.EndDate,
		Targeting: &request.Targeting{
			Gender:   campaign.Targeting.Gender,
			AgeFrom:  campaign.Targeting.AgeFrom,
			AgeTo:    ageTo,
			Location: campaign.Targeting.Location,
		},
	}

	updatedCampaign, err := h.advertisingClient.UpdateCampaign(updateRequestBody, user.PlatformID, campaignID)
	switch {
	case errors.Is(err, errorz.ErrInvalidData) || errors.Is(err, errorz.ErrCampaignNotFound):
		return c.Send(
			h.banners.Advertiser.Caption(h.layout.Text(c, "invalid_targeting_age_to")),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	case err != nil:
		return c.Send(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "campaign:update:back", struct {
				ID   string
				Page int
			}{
				ID:   campaignID,
				Page: p,
			}),
		)
	}

	return c.Send(
		h.banners.Advertiser.Caption(h.layout.Text(c, "campaign_updated_successfully", struct {
			AdTitle string
		}{
			AdTitle: updatedCampaign.AdTitle,
		})),
		h.layout.Markup(c, "campaign:update:back", struct {
			ID   string
			Page int
		}{
			ID:   campaignID,
			Page: p,
		}),
	)
}
