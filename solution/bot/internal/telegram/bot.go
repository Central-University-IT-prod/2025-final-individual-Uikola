package telegram

import (
	"bot/internal/telegram/advertiser"
	"bot/internal/telegram/client"
	"bot/internal/telegram/core"
	"bot/internal/telegram/menu"
	"bot/internal/telegram/middlewares"
	"bot/internal/telegram/start"

	"github.com/nlypage/intele"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/layout"
	"gopkg.in/telebot.v3/middleware"
)

func Setup(
	bot *tele.Bot,
	layout *layout.Layout,
	inputManager *intele.InputManager,

	coreHandler *core.Handler,
	middlewaresHandler *middlewares.Handler,
	startHandler *start.Handler,
	menuHandler *menu.Handler,
	clientHandler *client.Handler,
	advertiserHandler *advertiser.Handler,
) {

	// Pre-setup
	bot.Use(layout.Middleware("ru"))
	bot.Use(middleware.AutoRespond())
	bot.Handle(tele.OnText, inputManager.MessageHandler())
	bot.Handle(tele.OnMedia, inputManager.MessageHandler())
	bot.Handle(tele.OnCallback, inputManager.CallbackHandler())
	bot.Use(middlewaresHandler.ResetInputOnBack)

	bot.Handle(layout.Callback("core:hide"), coreHandler.Hide)
	bot.Handle(layout.Callback("core:cancel"), coreHandler.Hide)
	bot.Handle(layout.Callback("core:back"), coreHandler.Hide)

	// Handlers setup
	// Start
	bot.Handle("/start", startHandler.Start)

	// Auth
	bot.Handle(layout.Callback("auth:client"), clientHandler.ClientAuth)
	bot.Handle(layout.Callback("auth:advertiser"), advertiserHandler.AdvertiserAuth)
	bot.Handle(layout.Callback("auth:back_to_menu"), startHandler.BackToStart)
	bot.Use(middlewaresHandler.Authorized)

	// Main menu
	bot.Handle(layout.Callback("main_menu:logout"), startHandler.Logout)
	bot.Handle(layout.Callback("main_menu:back"), menuHandler.EditMenu)

	// Client
	clientHandler.ClientSetup(bot.Group(), middlewaresHandler)

	// Advertiser
	advertiserHandler.AdvertiserSetup(bot.Group(), middlewaresHandler)
}
