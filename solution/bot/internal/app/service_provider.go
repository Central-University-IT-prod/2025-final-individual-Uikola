package app

import (
	"bot/internal/config"
	"bot/internal/repository"
	"bot/internal/telegram/advertiser"
	"bot/internal/telegram/banners"
	"bot/internal/telegram/client"
	"bot/internal/telegram/core"
	"bot/internal/telegram/menu"
	"bot/internal/telegram/middlewares"
	"bot/internal/telegram/start"
	"bot/pkg/advertising"
	"bot/pkg/advertising/api"
	"context"
	"fmt"

	"github.com/nlypage/intele"

	"github.com/redis/go-redis/v9"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/layout"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	postgresRepository "bot/internal/repository/postgres"
	redisRepository "bot/internal/repository/redis"
)

// serviceProvider is a dependency container that lazily initializes dependencies.
type serviceProvider struct {
	telegramConfig    config.TelegramConfig
	loggerConfig      config.LoggerConfig
	pgConfig          config.PGConfig
	redisConfig       config.RedisConfig
	advertisingConfig config.AdvertisingConfig

	layout       *layout.Layout
	bot          *tele.Bot
	inputManager *intele.InputManager
	banners      *banners.Banners

	db                    *gorm.DB
	campaignRedisClient   *redis.Client
	stateRedisClient      *redis.Client
	advertiserRedisClient *redis.Client

	advertisingClient advertising.Client

	userRepository       repository.UserRepository
	stateRepository      repository.StateRepository
	campaignRepository   repository.CampaignRepository
	advertiserRepository repository.AdvertiserRepository

	coreHandler        *core.Handler
	middlewaresHandler *middlewares.Handler
	startHandler       *start.Handler
	menuHandler        *menu.Handler
	clientHandler      *client.Handler
	advertiserHandler  *advertiser.Handler
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

// TelegramConfig returns the telegram configuration.
func (s *serviceProvider) TelegramConfig() config.TelegramConfig {
	if s.telegramConfig == nil {
		cfg, err := config.NewTelegramConfig()
		if err != nil {
			panic(fmt.Errorf("failed to get telegram config: %w", err))
		}
		s.telegramConfig = cfg
	}

	return s.telegramConfig
}

// LoggerConfig returns the logger configuration.
func (s *serviceProvider) LoggerConfig() config.LoggerConfig {
	if s.loggerConfig == nil {
		cfg, err := config.NewLoggerConfig()
		if err != nil {
			panic(fmt.Errorf("failed to get logger config: %w", err))
		}

		s.loggerConfig = cfg
	}

	return s.loggerConfig
}

// PGConfig returns the postgres db configuration.
func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			panic(fmt.Errorf("failed to get pg config: %w", err))
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

// RedisConfig returns the redis configuration.
func (s *serviceProvider) RedisConfig() config.RedisConfig {
	if s.redisConfig == nil {
		cfg, err := config.NewRedisConfig()
		if err != nil {
			panic(fmt.Errorf("failed to get redis config: %w", err))
		}

		s.redisConfig = cfg
	}

	return s.redisConfig
}

// AdvertisingConfig returns the advertising configuration.
func (s *serviceProvider) AdvertisingConfig() config.AdvertisingConfig {
	if s.advertisingConfig == nil {
		cfg, err := config.NewAdvertisingConfig()
		if err != nil {
			panic(fmt.Errorf("failed to get advertising config: %w", err))
		}

		s.advertisingConfig = cfg
	}

	return s.advertisingConfig
}

// Layout returns the telebot layout.
func (s *serviceProvider) Layout() *layout.Layout {
	if s.layout == nil {
		lt, err := layout.New("telegram.yml")
		if err != nil {
			panic(fmt.Errorf("failed to get layout: %w", err))
		}

		s.layout = lt
	}

	return s.layout
}

// Bot returns the telebot bot.
func (s *serviceProvider) Bot() *tele.Bot {
	if s.bot == nil {
		b, err := tele.NewBot(s.Layout().Settings())
		if err != nil {
			panic(fmt.Errorf("failed to get bot: %w", err))
		}

		s.bot = b
	}

	return s.bot
}

// InputManager returns the intele input manager.
func (s *serviceProvider) InputManager(ctx context.Context) *intele.InputManager {
	if s.inputManager == nil {
		s.inputManager = intele.NewInputManager(intele.InputOptions{
			Storage: s.StateRepository(ctx),
		})
	}

	return s.inputManager
}

// Banners returns the telegram bot banners.
func (s *serviceProvider) Banners() *banners.Banners {
	if s.banners == nil {
		b, err := banners.New(s.Bot(), s.TelegramConfig().BotAuthBanner(), s.TelegramConfig().BotClientBanner(), s.TelegramConfig().BotAdvertiserBanner())
		if err != nil {
			panic(fmt.Errorf("failed to create banners: %w", err))
		}

		s.banners = b
	}

	return s.banners
}

// DBClient returns pg database connection.
func (s *serviceProvider) DBClient() *gorm.DB {
	if s.db == nil {
		gormConfig := &gorm.Config{
			TranslateError: true,
		}
		database, err := gorm.Open(postgres.Open(s.PGConfig().DSN()), gormConfig)
		if err != nil {
			panic(fmt.Errorf("failed to connect to database: %w", err))
		}

		s.db = database
	}

	return s.db
}

// CampaignRedisClient returns redis client for campaign repository.
func (s *serviceProvider) CampaignRedisClient(ctx context.Context) *redis.Client {
	if s.campaignRedisClient == nil {
		redisClient := redis.NewClient(&redis.Options{
			Addr:     s.RedisConfig().Addr(),
			Password: s.RedisConfig().Pass(),
			DB:       1,
		})
		if err := redisClient.Ping(ctx).Err(); err != nil {
			panic(fmt.Errorf("failed to ping state storage: %w", err))
		}

		s.campaignRedisClient = redisClient
	}

	return s.campaignRedisClient
}

// StateRedisClient returns redis client for state repository.
func (s *serviceProvider) StateRedisClient(ctx context.Context) *redis.Client {
	if s.stateRedisClient == nil {
		redisClient := redis.NewClient(&redis.Options{
			Addr:     s.RedisConfig().Addr(),
			Password: s.RedisConfig().Pass(),
			DB:       2,
		})
		if err := redisClient.Ping(ctx).Err(); err != nil {
			panic(fmt.Errorf("failed to ping state storage: %w", err))
		}

		s.stateRedisClient = redisClient
	}

	return s.stateRedisClient
}

// AdvertiserRedisClient returns redis client for state repository.
func (s *serviceProvider) AdvertiserRedisClient(ctx context.Context) *redis.Client {
	if s.advertiserRedisClient == nil {
		redisClient := redis.NewClient(&redis.Options{
			Addr:     s.RedisConfig().Addr(),
			Password: s.RedisConfig().Pass(),
			DB:       2,
		})
		if err := redisClient.Ping(ctx).Err(); err != nil {
			panic(fmt.Errorf("failed to ping state storage: %w", err))
		}

		s.advertiserRedisClient = redisClient
	}

	return s.advertiserRedisClient
}

// AdvertisingClient returns advertising platform client.
func (s *serviceProvider) AdvertisingClient() advertising.Client {
	if s.advertisingClient == nil {
		s.advertisingClient = api.NewClient(s.AdvertisingConfig().URL())
	}

	return s.advertisingClient
}

// UserRepository returns user repository.
func (s *serviceProvider) UserRepository() repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = postgresRepository.NewUserRepository(s.DBClient())
	}

	return s.userRepository
}

// StateRepository returns state repository.
func (s *serviceProvider) StateRepository(ctx context.Context) repository.StateRepository {
	if s.stateRepository == nil {
		s.stateRepository = redisRepository.NewStateRepository(s.StateRedisClient(ctx))
	}

	return s.stateRepository
}

// CampaignRepository returns state repository.
func (s *serviceProvider) CampaignRepository(ctx context.Context) repository.CampaignRepository {
	if s.campaignRepository == nil {
		s.campaignRepository = redisRepository.NewCampaignRepository(s.CampaignRedisClient(ctx))
	}

	return s.campaignRepository
}

// AdvertiserRepository returns advertiser repository.
func (s *serviceProvider) AdvertiserRepository(ctx context.Context) repository.AdvertiserRepository {
	if s.advertiserRepository == nil {
		s.advertiserRepository = redisRepository.NewAdvertiserRepository(s.AdvertiserRedisClient(ctx))
	}

	return s.advertiserRepository
}

// CoreHandler returns core handler.
func (s *serviceProvider) CoreHandler() *core.Handler {
	if s.coreHandler == nil {
		s.coreHandler = core.NewHandler()
	}

	return s.coreHandler
}

// MiddlewaresHandler returns middlewares handler.
func (s *serviceProvider) MiddlewaresHandler(ctx context.Context) *middlewares.Handler {
	if s.middlewaresHandler == nil {
		s.middlewaresHandler = middlewares.NewHandler(s.Bot(), s.Layout(), s.InputManager(ctx), s.Banners(), s.UserRepository())
	}

	return s.middlewaresHandler
}

// StartHandler returns start handler.
func (s *serviceProvider) StartHandler() *start.Handler {
	if s.startHandler == nil {
		s.startHandler = start.NewHandler(s.Layout(), s.Banners(), s.UserRepository(), s.MenuHandler())
	}

	return s.startHandler
}

// MenuHandler returns menu handler.
func (s *serviceProvider) MenuHandler() *menu.Handler {
	if s.menuHandler == nil {
		s.menuHandler = menu.NewHandler(s.Layout(), s.Banners(), s.UserRepository(), s.AdvertisingClient())
	}

	return s.menuHandler
}

// ClientHandler returns client handler.
func (s *serviceProvider) ClientHandler(ctx context.Context) *client.Handler {
	if s.clientHandler == nil {
		s.clientHandler = client.NewHandler(s.Layout(), s.InputManager(ctx), s.Banners(), s.UserRepository(), s.AdvertiserRepository(ctx), s.MenuHandler(), s.AdvertisingClient())
	}

	return s.clientHandler
}

// AdvertiserHandler returns advertiser handler.
func (s *serviceProvider) AdvertiserHandler(ctx context.Context) *advertiser.Handler {
	if s.advertiserHandler == nil {
		s.advertiserHandler = advertiser.NewHandler(s.Layout(), s.InputManager(ctx), s.Banners(), s.UserRepository(), s.CampaignRepository(ctx), s.MenuHandler(), s.AdvertisingClient())
	}

	return s.advertiserHandler
}
