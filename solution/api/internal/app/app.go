package app

import (
	"api/internal/closer"
	"api/internal/repository/clickhouse"
	"api/internal/repository/postgres"
	"api/internal/server"
	"api/pkg/zlog"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
)

// App represents the main application structure.
type App struct {
	serviceProvider *serviceProvider
	httpServer      *http.Server
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
func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	log.Info().Msg("Starting http server")
	return a.runHTTPServer()
}

// initializes application dependencies
func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initServiceProvider,
		a.initLogger,
		a.migratePG,
		a.migrateClickHouse,
		a.initHTTPServer,
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

func (a *App) migrateClickHouse(_ context.Context) error {
	return a.serviceProvider.ClickHouseClient().AutoMigrate(clickhouse.Migrations...)
}

func (a *App) initHTTPServer(ctx context.Context) error {
	srv := server.NewServer(
		a.serviceProvider.ClientHandler(),
		a.serviceProvider.AdvertiserHandler(),
		a.serviceProvider.MLScoreHandler(),
		a.serviceProvider.CampaignHandler(ctx),
		a.serviceProvider.TimeHandler(ctx),
		a.serviceProvider.StatisticHandler(),
	)

	httpServer := &http.Server{
		Addr:    a.serviceProvider.HTTPConfig().Address(),
		Handler: srv,
	}
	a.httpServer = httpServer

	return nil
}

func (a *App) runHTTPServer() error {
	if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().Err(err).Msg("Error starting http server")
		os.Exit(1)
	}
	return nil
}
