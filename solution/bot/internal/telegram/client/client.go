package client

import (
	"bot/internal/repository"
	"bot/internal/telegram/banners"
	"bot/internal/telegram/menu"
	"bot/internal/telegram/middlewares"
	"bot/pkg/advertising"

	"github.com/nlypage/intele"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/layout"
)

type Handler struct {
	layout       *layout.Layout
	inputManager *intele.InputManager
	banners      *banners.Banners

	userRepository       repository.UserRepository
	advertiserRepository repository.AdvertiserRepository

	menuHandler *menu.Handler

	advertisingClient advertising.Client
}

func NewHandler(layout *layout.Layout, inputManager *intele.InputManager, banners *banners.Banners, userRepository repository.UserRepository, advertiserRepository repository.AdvertiserRepository, menuHandler *menu.Handler, advertisingClient advertising.Client) *Handler {
	return &Handler{
		layout:       layout,
		inputManager: inputManager,
		banners:      banners,

		userRepository:       userRepository,
		advertiserRepository: advertiserRepository,

		menuHandler: menuHandler,

		advertisingClient: advertisingClient,
	}
}

func (h *Handler) ClientSetup(group *tele.Group, middlewaresHandler *middlewares.Handler) {
	group.Use(middlewaresHandler.IsClient)

	group.Handle(h.layout.Callback("main_menu:ads"), h.getAd)
	group.Handle(h.layout.Callback("ad:next"), h.getAd)
	group.Handle(h.layout.Callback("ad:click"), h.clickAd)
}
