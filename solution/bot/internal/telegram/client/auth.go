package client

import (
	"bot/internal/entity"
	"bot/internal/telegram/helper/validator"
	"bot/pkg/advertising/errorz"
	"context"
	"errors"

	"github.com/nlypage/intele/collector"
	"github.com/rs/zerolog/log"
	tele "gopkg.in/telebot.v3"
)

func (h *Handler) ClientAuth(c tele.Context) error {
	log.Info().Msgf("(user: %d) client auth", c.Sender().ID)
	inputCollector := collector.New()
	_ = c.Edit(
		h.banners.Auth.Caption(h.layout.Text(c, "input_client_id")),
		h.layout.Markup(c, "auth:back_to_menu"),
	)
	inputCollector.Collect(c.Message())

	var (
		clientID string
		done     bool
	)
	for {
		response, err := h.inputManager.Get(context.Background(), c.Sender().ID, 0)
		if response.Message != nil {
			inputCollector.Collect(response.Message)
		}
		switch {
		case response.Canceled:
			_ = inputCollector.Clear(c, collector.ClearOptions{IgnoreErrors: true, ExcludeLast: true})
			return nil
		case err != nil:
			log.Error().Msgf("(user: %d) error while input client id: %v", c.Sender().ID, err)
			_ = inputCollector.Send(c,
				h.banners.Auth.Caption(h.layout.Text(c, "input_error", h.layout.Text(c, "input_client_id"))),
				h.layout.Markup(c, "auth:back_to_menu"),
			)
		case response.Message == nil:
			log.Error().Msgf("(user: %d) error while input client id: %v", c.Sender().ID, err)
			_ = inputCollector.Send(c,
				h.banners.Auth.Caption(h.layout.Text(c, "input_error", h.layout.Text(c, "input_client_id"))),
				h.layout.Markup(c, "auth:back_to_menu"),
			)
		case !validator.ClientID(response.Message.Text):
			_ = inputCollector.Send(c,
				h.banners.Auth.Caption(h.layout.Text(c, "invalid_client_id")),
				h.layout.Markup(c, "auth:back_to_menu"),
			)
		case validator.ClientID(response.Message.Text):
			clientID = response.Message.Text
			_, err = h.advertisingClient.GetClientByID(clientID)
			switch {
			case errors.Is(err, errorz.ErrInvalidData) || errors.Is(err, errorz.ErrClientNotFound):
				_ = inputCollector.Send(c,
					h.banners.Auth.Caption(h.layout.Text(c, "invalid_client_id")),
					h.layout.Markup(c, "auth:back_to_menu"),
				)
				continue
			case err != nil:
				_ = inputCollector.Send(c,
					h.banners.Auth.Caption(h.layout.Text(c, "technical_issues", err.Error())),
					h.layout.Markup(c, "auth:back_to_menu"),
				)
				continue
			}

			_ = inputCollector.Clear(c, collector.ClearOptions{IgnoreErrors: true})
			done = true
		}
		if done {
			break
		}
	}

	user := entity.User{
		ID:         c.Sender().ID,
		Username:   c.Sender().Username,
		PlatformID: clientID,
		Role:       entity.Client,
	}
	_, err := h.userRepository.Save(context.Background(), user)
	if err != nil {
		log.Error().Msgf("(user: %d) error while saving user: %v", c.Sender().ID, err)
		return c.Send(
			h.banners.Auth.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "auth:back_to_menu"),
		)
	}
	log.Info().Msgf("(user: %d) user saved(role: %s)", c.Sender().ID, user.Role)

	return h.menuHandler.SendMenu(c)
}
