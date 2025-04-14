package app

import (
	"bot/internal/closer"
	"bot/internal/repository/postgres"
	"bot/internal/telegram"
	"bot/pkg/zlog"
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	tele "gopkg.in/telebot.v3"
)

// App represents the main application structure.
type App struct {
	serviceProvider *serviceProvider
	bot             *tele.Bot
}

// NewApp initializes the application and its dependencies.
func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, fmt.Errorf("new app: %w", err)
	}

	return a, nil
}

// Run starts the application.
func (a *App) Run() {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	log.Info().Msg("Starting bot")
	a.bot.Start()
}

// initializes application dependencies
func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initServiceProvider,
		a.initLogger,
		a.migratePG,
		a.initBot,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return fmt.Errorf("init deps: %w", err)
		}
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	cfg := a.serviceProvider.LoggerConfig()

	log.Logger = zlog.Default(cfg.IsPretty(), cfg.Version(), cfg.LogLevel())

	return nil
}

func (a *App) migratePG(_ context.Context) error {
	return a.serviceProvider.DBClient().AutoMigrate(postgres.Migrations...)
}

func (a *App) initBot(ctx context.Context) error {
	if cmds := a.serviceProvider.Layout().Commands(); cmds != nil {
		if err := a.serviceProvider.Bot().SetCommands(cmds); err != nil {
			return err
		}
	}

	telegram.Setup(
		a.serviceProvider.Bot(),
		a.serviceProvider.Layout(),
		a.serviceProvider.InputManager(ctx),
		a.serviceProvider.CoreHandler(),
		a.serviceProvider.MiddlewaresHandler(ctx),
		a.serviceProvider.StartHandler(),
		a.serviceProvider.MenuHandler(),
		a.serviceProvider.ClientHandler(ctx),
		a.serviceProvider.AdvertiserHandler(ctx),
	)
	a.bot = a.serviceProvider.Bot()

	return nil
}
