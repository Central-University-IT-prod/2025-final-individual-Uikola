package client

import (
	"bot/pkg/advertising/errorz"
	"bot/pkg/advertising/request"
	"bot/pkg/advertising/response"
	"context"
	"errors"

	tele "gopkg.in/telebot.v3"
)

func (h *Handler) getAd(c tele.Context) error {
	h.advertiserRepository.Clear(c.Sender().ID)

	user, err := h.userRepository.GetByID(context.Background(), c.Sender().ID)
	if err != nil {
		return c.Edit(
			h.banners.Client.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "main_menu:back"),
		)
	}

	client, err := h.advertisingClient.GetClientByID(user.PlatformID)
	if err != nil {
		return c.Edit(
			h.banners.Client.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "main_menu:back"),
		)
	}

	ad, err := h.advertisingClient.GetAd(client.ClientID)
	switch {
	case errors.Is(err, errorz.ErrCampaignNotFound):
		return c.Edit(
			h.banners.Client.Caption(h.layout.Text(c, "no_ads_found")),
			h.layout.Markup(c, "main_menu:back"),
		)
	case err != nil:
		return c.Edit(
			h.banners.Client.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "main_menu:back"),
		)
	}

	advertiser, err := h.advertisingClient.GetAdvertiserByID(ad.AdvertiserID)
	if err != nil {
		return c.Edit(
			h.banners.Client.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "main_menu:back"),
		)
	}

	h.advertiserRepository.Set(c.Sender().ID, advertiser, 0)

	return c.Edit(
		h.banners.Client.Caption(h.layout.Text(c, "ad", struct {
			response.GetAd
			AdvertiserName string
		}{
			GetAd:          ad,
			AdvertiserName: advertiser.Name,
		})),
		h.layout.Markup(c, "ad:menu", struct {
			AdID string
		}{
			AdID: ad.AdID,
		}),
	)
}

func (h *Handler) clickAd(c tele.Context) error {
	campaignID := c.Callback().Data

	advertiser, err := h.advertiserRepository.Get(c.Sender().ID)
	if err != nil {
		return c.Edit(
			h.banners.Client.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "main_menu:back"),
		)
	}

	user, err := h.userRepository.GetByID(context.Background(), c.Sender().ID)
	if err != nil {
		return c.Edit(
			h.banners.Client.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "main_menu:back"),
		)
	}

	req := request.ClickAd{
		ClientID: user.PlatformID,
	}
	err = h.advertisingClient.ClickAd(req, campaignID)
	switch {
	case errors.Is(err, errorz.ErrImpressionNotFound):
		return c.Edit(
			h.banners.Client.Caption(h.layout.Text(c, "not_seen_ad")),
			h.layout.Markup(c, "main_menu:back"),
		)
	case err != nil:
		return c.Edit(
			h.banners.Client.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "main_menu:back"),
		)
	}

	campaign, err := h.advertisingClient.GetCampaign(advertiser.AdvertiserID, campaignID)
	if err != nil {
		return c.Edit(
			h.banners.Client.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "main_menu:back"),
		)
	}

	return c.Edit(
		h.banners.Client.Caption(h.layout.Text(c, "ad_with_details", struct {
			AdTitle        string
			AdText         string
			AdvertiserName string
			AdID           string
			ImageURL       string
		}{
			AdTitle:        campaign.AdTitle,
			AdText:         campaign.AdText,
			AdvertiserName: advertiser.Name,
			AdID:           campaign.CampaignID,
			ImageURL:       campaign.ImageURL,
		})),
		h.layout.Markup(c, "clicked_ad:menu"),
	)
}
