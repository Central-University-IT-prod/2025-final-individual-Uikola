package menu

import (
	"bot/internal/entity"
	"bot/internal/repository"
	"bot/internal/telegram/banners"
	"bot/pkg/advertising"
	"context"

	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/layout"
)

type Handler struct {
	layout  *layout.Layout
	banners *banners.Banners

	userRepository repository.UserRepository

	advertisingClient advertising.Client
}

func NewHandler(layout *layout.Layout, banners *banners.Banners, userRepository repository.UserRepository, advertisingClient advertising.Client) *Handler {
	return &Handler{
		layout:  layout,
		banners: banners,

		userRepository: userRepository,

		advertisingClient: advertisingClient,
	}
}

func (h *Handler) SendMenu(c tele.Context) error {
	user, err := h.userRepository.GetByID(context.Background(), c.Sender().ID)
	if err != nil {
		return c.Send(
			h.banners.Auth.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "core:hide"),
		)
	}

	currentDate, err := h.advertisingClient.GetCurrentDate()
	if err != nil {
		return c.Send(
			h.banners.Auth.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "core:hide"),
		)
	}

	switch user.Role {
	case entity.Client:
		return c.Send(
			h.banners.Client.Caption(h.layout.Text(c, "main_menu_text", struct {
				CurrentDate int
			}{
				CurrentDate: currentDate,
			})),
			h.layout.Markup(c, "main_menu:client:menu"),
		)
	case entity.Advertiser:
		return c.Send(
			h.banners.Advertiser.Caption(h.layout.Text(c, "main_menu_text", struct {
				CurrentDate int
			}{
				CurrentDate: currentDate,
			})),
			h.layout.Markup(c, "main_menu:advertiser:menu"),
		)
	default:
		return c.Send(
			h.banners.Auth.Caption(h.layout.Text(c, "something_went_wrong")),
			h.layout.Markup(c, "core:hide"),
		)
	}
}

func (h *Handler) EditMenu(c tele.Context) error {
	user, err := h.userRepository.GetByID(context.Background(), c.Sender().ID)
	if err != nil {
		return c.Send(
			h.banners.Auth.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "core:hide"),
		)
	}

	currentDate, err := h.advertisingClient.GetCurrentDate()
	if err != nil {
		return c.Send(
			h.banners.Auth.Caption(h.layout.Text(c, "technical_issues", err.Error())),
			h.layout.Markup(c, "core:hide"),
		)
	}

	switch user.Role {
	case entity.Client:
		return c.Edit(
			h.banners.Client.Caption(h.layout.Text(c, "main_menu_text", struct {
				CurrentDate int
			}{
				CurrentDate: currentDate,
			})),
			h.layout.Markup(c, "main_menu:client:menu"),
		)
	case entity.Advertiser:
		return c.Edit(
			h.banners.Advertiser.Caption(h.layout.Text(c, "main_menu_text", struct {
				CurrentDate int
			}{
				CurrentDate: currentDate,
			})),
			h.layout.Markup(c, "main_menu:advertiser:menu"),
		)
	default:
		return c.Edit(
			h.banners.Auth.Caption(h.layout.Text(c, "something_went_wrong")),
			h.layout.Markup(c, "core:hide"),
		)
	}
}
