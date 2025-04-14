package advertiser

import (
	"bot/internal/telegram/helper/validator"
	"bot/pkg/advertising/request"
	"context"

	"github.com/nlypage/intele/collector"
	"github.com/rs/zerolog/log"
	tele "gopkg.in/telebot.v3"
)

func (h *Handler) generateAdText(c tele.Context) error {
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

	inputCollector := collector.New()
	inputCollector.Collect(c.Message())

	steps := []struct {
		promptKey   string
		errorKey    string
		result      *string
		validator   func(string, map[string]interface{}) bool
		callbackBtn *tele.Btn
	}{
		{
			promptKey:   "input_ad_title_for_generation",
			errorKey:    "invalid_ad_title_for_generation",
			result:      new(string),
			validator:   validator.AdTitle,
			callbackBtn: nil,
		},

		{
			promptKey:   "input_context_for_generation",
			errorKey:    "invalid_context_for_generation",
			result:      new(string),
			validator:   validator.ContextForGeneration,
			callbackBtn: h.layout.Button(c, "generate_ad_text:context_skip"),
		},
	}

	isFirst := true
	for _, step := range steps {
		done := false

		markup := h.layout.Markup(c, "main_menu:back")
		if step.callbackBtn != nil {
			markup.InlineKeyboard = append(
				[][]tele.InlineButton{{*step.callbackBtn.Inline()}},
				markup.InlineKeyboard...,
			)
		}

		if isFirst {
			_ = c.Edit(
				h.banners.Advertiser.Caption(h.layout.Text(c, step.promptKey)),
				markup,
			)
		} else {
			_ = inputCollector.Send(c,
				h.banners.Advertiser.Caption(h.layout.Text(c, step.promptKey)),
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
					h.banners.Advertiser.Caption(h.layout.Text(c, "input_error", h.layout.Text(c, step.promptKey))),
					h.layout.Markup(c, "main_menu:back"),
				)
			case resp.Callback != nil:
				*step.result = ""
				_ = inputCollector.Clear(c, collector.ClearOptions{IgnoreErrors: true})
				done = true
			case !step.validator(resp.Message.Text, nil):
				_ = inputCollector.Send(c,
					h.banners.Advertiser.Caption(h.layout.Text(c, step.errorKey)),
					h.layout.Markup(c, "main_menu:back"),
				)
			case step.validator(resp.Message.Text, nil):
				*step.result = resp.Message.Text
				_ = inputCollector.Clear(c, collector.ClearOptions{IgnoreErrors: true})
				done = true
			}
		}
	}

	adTitle := *steps[0].result
	ctx := *steps[1].result

	req := request.GenerateAdText{
		AdTitle: adTitle,
		Context: ctx,
	}

	resp, err := h.advertisingClient.GenerateAdText(req, advertiser.AdvertiserID)
	if err != nil {
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "main_menu:back"),
		)
	}

	return c.Send(
		h.banners.Advertiser.Caption(h.layout.Text(c, "text_generated_successfully", resp)),
		h.layout.Markup(c, "main_menu:back"),
	)
}
