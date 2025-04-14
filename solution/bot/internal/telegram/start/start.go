package start

import (
	"bot/internal/errorz"
	"bot/internal/repository"
	"bot/internal/telegram/banners"
	"bot/internal/telegram/menu"
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/layout"
)

type Handler struct {
	layout  *layout.Layout
	banners *banners.Banners

	userRepository repository.UserRepository

	menuHandler *menu.Handler
}

func NewHandler(layout *layout.Layout, banners *banners.Banners, userRepository repository.UserRepository, menuHandler *menu.Handler) *Handler {
	return &Handler{
		layout:  layout,
		banners: banners,

		userRepository: userRepository,

		menuHandler: menuHandler,
	}
}

func (h *Handler) Start(c tele.Context) error {
	log.Info().Msgf("(user: %d) enter /start", c.Sender().ID)

	_, err := h.userRepository.GetByID(context.Background(), c.Sender().ID)
	if err != nil && !errors.Is(err, errorz.ErrUserNotFound) {
		log.Error().Msgf("(user: %d) error while getting user from db: %v", c.Sender().ID, err)
		return c.Send(
			h.banners.Auth.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "core:hide"),
		)
	}

	if errors.Is(err, errorz.ErrUserNotFound) {
		return c.Send(
			h.banners.Auth.Caption(h.layout.Text(c, "auth_menu_text", struct {
				Username string
			}{
				Username: c.Sender().Username,
			})),
			h.layout.Markup(c, "auth:menu"),
		)
	}

	return h.menuHandler.SendMenu(c)
}

func (h *Handler) BackToStart(c tele.Context) error {
	log.Info().Msgf("(user: %d) back to start", c.Sender().ID)

	_, err := h.userRepository.GetByID(context.Background(), c.Sender().ID)
	if err != nil && !errors.Is(err, errorz.ErrUserNotFound) {
		log.Error().Msgf("(user: %d) error while getting user from db: %v", c.Sender().ID, err)
		return c.Send(
			h.banners.Auth.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "core:hide"),
		)
	}

	if errors.Is(err, errorz.ErrUserNotFound) {
		return c.Edit(
			h.banners.Auth.Caption(h.layout.Text(c, "auth_menu_text", struct {
				Username string
			}{
				Username: c.Sender().Username,
			})),
			h.layout.Markup(c, "auth:menu"),
		)
	}

	return h.menuHandler.EditMenu(c)
}

func (h *Handler) Logout(c tele.Context) error {
	log.Info().Msgf("(user: %d) logout", c.Sender().ID)

	err := h.userRepository.Delete(context.Background(), c.Sender().ID)
	if err != nil {
		log.Error().Msgf("(user: %d) error deleting user: %v", c.Sender().ID, err)
		return c.Send(
			h.banners.Auth.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "core:hide"),
		)
	}

	_ = c.Delete()

	return c.Send(
		h.banners.Auth.Caption(h.layout.Text(c, "logout_success")),
		h.layout.Markup(c, "core:hide"),
	)
}
