package middlewares

import (
	"bot/internal/entity"
	"bot/internal/errorz"
	"bot/internal/repository"
	"bot/internal/telegram/banners"
	"context"
	"errors"
	"strings"

	"github.com/nlypage/intele"

	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/layout"
)

type Handler struct {
	bot          *tele.Bot
	layout       *layout.Layout
	inputManager *intele.InputManager
	banners      *banners.Banners

	userRepository repository.UserRepository
}

func NewHandler(bot *tele.Bot, layout *layout.Layout, inputManager *intele.InputManager, banners *banners.Banners, userRepository repository.UserRepository) *Handler {
	return &Handler{
		bot:          bot,
		layout:       layout,
		inputManager: inputManager,
		banners:      banners,

		userRepository: userRepository,
	}
}

func (h *Handler) Authorized(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		user, err := h.userRepository.GetByID(context.Background(), c.Sender().ID)
		if err != nil {
			if !errors.Is(err, errorz.ErrUserNotFound) {
				return c.Send(
					h.banners.Auth.Caption(h.layout.Text(c, "technical_issues", err.Error())),
					h.layout.Markup(c, "core:hide"),
				)
			}
			return c.Send(
				h.banners.Auth.Caption(h.layout.Text(c, "auth_required")),
				h.layout.Markup(c, "core:hide"),
			)
		}

		if c.Sender().Username != user.Username {
			user.Username = c.Sender().Username
			_, err = h.userRepository.Save(context.Background(), user)
			if err != nil {
				return c.Send(
					h.banners.Auth.Caption(h.layout.Text(c, "technical_issues", err.Error())),
					h.layout.Markup(c, "core:hide"),
				)
			}
		}

		return next(c)
	}
}

// ResetInputOnBack middleware clears the input state when the back button is pressed.
func (h *Handler) ResetInputOnBack(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if c.Callback() != nil {
			if strings.Contains(c.Callback().Data, "back") || strings.Contains(c.Callback().Unique, "back") {
				h.inputManager.Cancel(c.Sender().ID)
			}
		}
		if c.Message() != nil {
			if strings.HasPrefix(c.Message().Text, "/") {
				h.inputManager.Cancel(c.Sender().ID)
			}
		}

		return next(c)
	}
}

func (h *Handler) IsAdvertiser(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		user, err := h.userRepository.GetByID(context.Background(), c.Sender().ID)
		if err != nil {
			return c.Send(
				h.banners.Advertiser.Caption(h.layout.Text(c, "technical_issues", err.Error())),
				h.layout.Markup(c, "core:hide"),
			)
		}

		if user.Role != entity.Advertiser {
			return c.Send(
				h.banners.Advertiser.Caption(h.layout.Text(c, "advertiser_required")),
				h.layout.Markup(c, "core:hide"),
			)
		}

		return next(c)
	}
}

func (h *Handler) IsClient(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		user, err := h.userRepository.GetByID(context.Background(), c.Sender().ID)
		if err != nil {
			return c.Send(
				h.banners.Client.Caption(h.layout.Text(c, "technical_issues", err.Error())),
				h.layout.Markup(c, "core:hide"),
			)
		}

		if user.Role != entity.Client {
			return c.Send(
				h.banners.Client.Caption(h.layout.Text(c, "client_required")),
				h.layout.Markup(c, "core:hide"),
			)
		}

		return next(c)
	}
}
