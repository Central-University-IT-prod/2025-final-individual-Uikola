package advertiser

import (
	"bot/internal/telegram/helper/validator"
	"bot/pkg/advertising/errorz"
	"bot/pkg/advertising/request"
	"context"
	"errors"
	"strconv"

	"github.com/nlypage/intele/collector"
	"github.com/rs/zerolog/log"
	tele "gopkg.in/telebot.v3"
)

func (h *Handler) createCampaign(c tele.Context) error {
	log.Info().Msgf("(user: %d) create new capaing request(club=%s)", c.Sender().ID)

	h.campaignRepository.Clear(c.Sender().ID)

	_, err := h.userRepository.GetByID(context.Background(), c.Sender().ID)
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
			h.layout.Markup(c, "main_menu:back"),
		)
	}

	inputCollector := collector.New()
	inputCollector.Collect(c.Message())

	isFirst := true

	var steps []struct {
		promptKey   string
		objectFunc  func() interface{}
		errorKey    string
		result      *string
		validator   func(string, map[string]interface{}) bool
		paramsFunc  func(map[string]interface{}) map[string]interface{}
		callbackBtn *tele.Btn
	}

	steps = []struct {
		promptKey   string
		objectFunc  func() interface{}
		errorKey    string
		result      *string
		validator   func(string, map[string]interface{}) bool
		paramsFunc  func(map[string]interface{}) map[string]interface{}
		callbackBtn *tele.Btn
	}{
		{
			promptKey: "input_impressions_limit",
			objectFunc: func() interface{} {
				return struct{}{}
			},
			errorKey:    "invalid_impressions_limit",
			result:      new(string),
			validator:   validator.ImpressionsLimit,
			paramsFunc:  nil,
			callbackBtn: nil,
		},

		{
			promptKey: "input_clicks_limit",
			objectFunc: func() interface{} {
				return struct{}{}
			},
			errorKey:    "invalid_clicks_limit",
			result:      new(string),
			validator:   validator.ClicksLimit,
			paramsFunc:  nil,
			callbackBtn: nil,
		},

		{
			promptKey: "input_cost_per_impression",
			objectFunc: func() interface{} {
				return struct{}{}
			},
			errorKey:    "invalid_cost_per_impression",
			result:      new(string),
			validator:   validator.CostPerImpression,
			paramsFunc:  nil,
			callbackBtn: nil,
		},
		{
			promptKey: "input_cost_per_click",
			objectFunc: func() interface{} {
				return struct{}{}
			},
			errorKey:    "invalid_cost_per_click",
			result:      new(string),
			validator:   validator.CostPerClick,
			paramsFunc:  nil,
			callbackBtn: nil,
		},
		{
			promptKey: "input_ad_title",
			objectFunc: func() interface{} {
				return struct{}{}
			},
			errorKey:    "invalid_at_title",
			result:      new(string),
			validator:   validator.AdTitle,
			paramsFunc:  nil,
			callbackBtn: nil,
		},
		{
			promptKey: "input_ad_text",
			objectFunc: func() interface{} {
				return struct{}{}
			},
			errorKey:    "invalid_ad_text",
			result:      new(string),
			validator:   validator.AdText,
			paramsFunc:  nil,
			callbackBtn: nil,
		},
		{
			promptKey: "input_start_date",
			objectFunc: func() interface{} {
				return struct {
					CurrentDate int
				}{
					CurrentDate: currentDate,
				}
			},
			errorKey:  "invalid_start_date",
			result:    new(string),
			validator: validator.StartDate,
			paramsFunc: func(params map[string]interface{}) map[string]interface{} {
				if params == nil {
					params = make(map[string]interface{})
				}
				params["currentDate"] = currentDate
				return params
			},
			callbackBtn: nil,
		},
		{
			promptKey: "input_end_date",
			objectFunc: func() interface{} {
				return struct {
					StartDate string
				}{
					StartDate: *steps[6].result,
				}
			},
			errorKey:  "invalid_end_date",
			result:    new(string),
			validator: validator.EndDate,
			paramsFunc: func(params map[string]interface{}) map[string]interface{} {
				if params == nil {
					params = make(map[string]interface{})
				}
				params["startDate"] = *steps[6].result
				return params
			},
			callbackBtn: nil,
		},
		{
			promptKey: "input_targeting_gender",
			objectFunc: func() interface{} {
				return struct{}{}
			},
			errorKey:    "invalid_targeting_gender",
			result:      new(string),
			validator:   validator.TargetingGender,
			paramsFunc:  nil,
			callbackBtn: h.layout.Button(c, "advertiser:create_campaign:targeting_gender_skip"),
		},
		{
			promptKey: "input_targeting_age_from",
			objectFunc: func() interface{} {
				return struct{}{}
			},
			errorKey:    "invalid_targeting_age_from",
			result:      new(string),
			validator:   validator.TargetingAgeFrom,
			paramsFunc:  nil,
			callbackBtn: h.layout.Button(c, "advertiser:create_campaign:targeting_age_from_skip"),
		},
		{
			promptKey: "input_targeting_age_to",
			objectFunc: func() interface{} {
				return struct{}{}
			},
			errorKey:  "invalid_targeting_age_to",
			result:    new(string),
			validator: validator.TargetingAgeTo,
			paramsFunc: func(params map[string]interface{}) map[string]interface{} {
				if params == nil {
					params = make(map[string]interface{})
				}
				params["ageFrom"] = *steps[9].result
				return params
			},
			callbackBtn: h.layout.Button(c, "advertiser:create_campaign:targeting_age_to_skip"),
		},
		{
			promptKey: "input_targeting_location",
			objectFunc: func() interface{} {
				return struct{}{}
			},
			errorKey:    "invalid_targeting_location",
			result:      new(string),
			validator:   validator.TargetingLocation,
			paramsFunc:  nil,
			callbackBtn: h.layout.Button(c, "advertiser:create_campaign:targeting_location_skip"),
		},
	}

	for _, step := range steps {
		done := false

		var params map[string]interface{}
		if step.paramsFunc != nil {
			params = step.paramsFunc(params)
		}

		markup := h.layout.Markup(c, "main_menu:back")
		if step.callbackBtn != nil {
			markup.InlineKeyboard = append(
				[][]tele.InlineButton{{*step.callbackBtn.Inline()}},
				markup.InlineKeyboard...,
			)
		}

		if isFirst {
			_ = c.Edit(
				h.banners.Advertiser.Caption(h.layout.Text(c, step.promptKey, step.objectFunc())),
				markup,
			)
		} else {
			_ = inputCollector.Send(c,
				h.banners.Advertiser.Caption(h.layout.Text(c, step.promptKey, step.objectFunc())),
				markup,
			)
		}
		isFirst = false

		for !done {
			resp, errGet := h.inputManager.Get(context.Background(), c.Sender().ID, 0, step.callbackBtn)
			if resp.Message != nil {
				inputCollector.Collect(resp.Message)
			}
			switch {
			case resp.Canceled:
				_ = inputCollector.Clear(c, collector.ClearOptions{IgnoreErrors: true, ExcludeLast: true})
				return nil
			case errGet != nil:
				log.Error().Msgf("(user: %d) error while input step (%s): %v", c.Sender().ID, step.promptKey, errGet)
				_ = inputCollector.Send(c,
					h.banners.Advertiser.Caption(h.layout.Text(c, "input_error", h.layout.Text(c, step.promptKey, step.objectFunc()))),
					h.layout.Markup(c, "main_menu:back"),
				)
			case resp.Callback != nil:
				*step.result = ""
				_ = inputCollector.Clear(c, collector.ClearOptions{IgnoreErrors: true})
				done = true
			case !step.validator(resp.Message.Text, params):
				_ = inputCollector.Send(c,
					h.banners.Advertiser.Caption(h.layout.Text(c, step.errorKey, step.objectFunc())),
					h.layout.Markup(c, "main_menu:back"),
				)
			case step.validator(resp.Message.Text, params):
				*step.result = resp.Message.Text
				_ = inputCollector.Clear(c, collector.ClearOptions{IgnoreErrors: true})
				done = true
			}
		}
	}

	campaignImpressionsLimit, _ := strconv.Atoi(*steps[0].result)
	campaignClicksLimit, _ := strconv.Atoi(*steps[1].result)
	campaignCostPerImpression, _ := strconv.ParseFloat(*steps[2].result, 64)
	campaignCostPerClick, _ := strconv.ParseFloat(*steps[3].result, 64)
	campaignAdTitle := *steps[4].result
	campaignAdText := *steps[5].result
	campaignStartDate, _ := strconv.Atoi(*steps[6].result)
	campaignEndDate, _ := strconv.Atoi(*steps[7].result)

	var campaignTargetingGender *string
	campaignTargetingGenderStr := ""
	if *steps[8].result != "" {
		campaignTargetingGender = steps[8].result
		campaignTargetingGenderStr = *steps[8].result
	}

	var campaignTargetingAgeFrom, campaignTargetingAgeTo *int
	var campaignTargetingAgeFromInt, campaignTargetingAgeToInt int
	if *steps[9].result != "" {
		temp, _ := strconv.Atoi(*steps[9].result)
		campaignTargetingAgeFrom = &temp
		campaignTargetingAgeFromInt = temp
	}
	if *steps[10].result != "" {
		temp, _ := strconv.Atoi(*steps[10].result)
		campaignTargetingAgeTo = &temp
		campaignTargetingAgeToInt = temp
	}

	var campaignTargetingLocation *string
	campaignTargetingLocationStr := ""
	if *steps[11].result != "" {
		campaignTargetingLocation = steps[11].result
		campaignTargetingLocationStr = *steps[11].result
	}

	createCampaignRequestBody := request.CreateCampaign{
		ImpressionsLimit:  campaignImpressionsLimit,
		ClicksLimit:       campaignClicksLimit,
		CostPerImpression: campaignCostPerImpression,
		CostPerClick:      campaignCostPerClick,
		AdTitle:           campaignAdTitle,
		AdText:            campaignAdText,
		StartDate:         campaignStartDate,
		EndDate:           campaignEndDate,
		Targeting: request.Targeting{
			Gender:   campaignTargetingGender,
			AgeFrom:  campaignTargetingAgeFrom,
			AgeTo:    campaignTargetingAgeTo,
			Location: campaignTargetingLocation,
		},
	}

	h.campaignRepository.Set(c.Sender().ID, createCampaignRequestBody, 0)

	return c.Send(
		h.layout.Text(c, "campaign_confirmation", struct {
			ImpressionsLimit  int
			ClicksLimit       int
			CostPerImpression float64
			CostPerClick      float64
			AdTitle           string
			AdText            string
			StartDate         int
			EndDate           int
			Gender            string
			AgeFrom           int
			AgeTo             int
			Location          string
		}{
			ImpressionsLimit:  campaignImpressionsLimit,
			ClicksLimit:       campaignClicksLimit,
			CostPerImpression: campaignCostPerImpression,
			CostPerClick:      campaignCostPerClick,
			AdTitle:           campaignAdTitle,
			AdText:            campaignAdText,
			StartDate:         campaignStartDate,
			EndDate:           campaignEndDate,
			Gender:            campaignTargetingGenderStr,
			AgeFrom:           campaignTargetingAgeFromInt,
			AgeTo:             campaignTargetingAgeToInt,
			Location:          campaignTargetingLocationStr,
		}),
		h.layout.Markup(c, "advertiser:create_campaign:confirm"),
	)
}

func (h *Handler) confirmCampaignCreation(c tele.Context) error {
	user, err := h.userRepository.GetByID(context.Background(), c.Sender().ID)
	if err != nil {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "main_menu:back"),
		)
	}

	campaign, err := h.campaignRepository.Get(c.Sender().ID)
	if err != nil {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "main_menu:back"),
		)
	}

	createdCampaign, err := h.advertisingClient.CreateCampaign(campaign, user.PlatformID)
	switch {
	case errors.Is(err, errorz.ErrInvalidData) || errors.Is(err, errorz.ErrAdvertiserNotFound):
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "invalid_campaign_data")),
			h.layout.Markup(c, "main_menu:back"),
		)
	case err != nil:
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "main_menu:back"),
		)
	}

	return c.Edit(
		h.banners.Advertiser.Caption(h.layout.Text(c, "campaign_created", struct {
			AdTitle string
		}{
			AdTitle: createdCampaign.AdTitle,
		})),
		h.layout.Markup(c, "main_menu:back"),
	)
}
