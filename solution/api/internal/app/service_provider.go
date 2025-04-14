package app

import (
	"api/internal/config"
	"api/internal/repository"
	"api/internal/server/advertiser"
	"api/internal/server/campaign"
	"api/internal/server/client"
	"api/internal/server/mlscore"
	"api/internal/server/statistic"
	"api/internal/server/time"
	"api/internal/usecase"
	"api/pkg/ai"
	"api/pkg/ai/gigachat"
	"api/pkg/s3"
	"api/pkg/s3/minio"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"gorm.io/gorm"

	clickHouseRepository "api/internal/repository/clickhouse"
	postgresRepository "api/internal/repository/postgres"
	redisRepository "api/internal/repository/redis"
	advertiserUseCase "api/internal/usecase/advertiser"
	campaignUsecase "api/internal/usecase/campaign"
	clientUsecase "api/internal/usecase/client"
	mlScoreUsecase "api/internal/usecase/mlscore"
	statisticUsecase "api/internal/usecase/statistic"

	"gorm.io/driver/clickhouse"
	"gorm.io/driver/postgres"
)

// serviceProvider is a dependency container that lazily initializes dependencies.
type serviceProvider struct {
	httpConfig       config.HTTPConfig
	loggerConfig     config.LoggerConfig
	pgConfig         config.PGConfig
	redisConfig      config.RedisConfig
	clickHouseConfig config.ClickHouseConfig
	minioConfig      config.MinioConfig
	gigaChatConfig   config.GigaChatConfig
	moderationConfig config.ModerationConfig

	db               *gorm.DB
	redisClient      *redis.Client
	clickHouseClient *gorm.DB
	s3Client         s3.Client
	aiClient         ai.Client

	clientRepository     repository.ClientRepository
	advertiserRepository repository.AdvertiserRepository
	mlScoreRepository    repository.MLScoreRepository
	campaignRepository   repository.CampaignRepository
	timeRepository       repository.TimeRepository
	displayRepository    repository.ImpressionRepository
	clickRepository      repository.ClickRepository
	statisticRepository  repository.StatisticRepository
	adRepository         repository.AdRepository

	clientUsecase     usecase.ClientUsecase
	advertiserUsecase usecase.AdvertiserUsecase
	mlScoreUsecase    usecase.MLScoreUsecase
	campaignUsecase   usecase.CampaignUsecase
	statisticUsecase  usecase.StatisticUsecase

	clientHandler     *client.Handler
	advertiserHandler *advertiser.Handler
	mlScoreHandler    *mlscore.Handler
	campaignHandler   *campaign.Handler
	timeHandler       *time.Handler
	statisticHandler  *statistic.Handler
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

// HTTPConfig returns the HTTP configuration.
func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := config.NewHTTPConfig()
		if err != nil {
			panic(fmt.Errorf("failed to get http config: %w", err))
		}
		s.httpConfig = cfg
	}

	return s.httpConfig
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

// ClickHouseConfig returns the clickhouse configuration.
func (s *serviceProvider) ClickHouseConfig() config.ClickHouseConfig {
	if s.clickHouseConfig == nil {
		cfg, err := config.NewClickHouseConfig()
		if err != nil {
			panic(fmt.Errorf("failed to get clickhouse config: %w", err))
		}

		s.clickHouseConfig = cfg
	}

	return s.clickHouseConfig
}

// MinioConfig returns the minio configuration.
func (s *serviceProvider) MinioConfig() config.MinioConfig {
	if s.minioConfig == nil {
		cfg, err := config.NewMinioConfig()
		if err != nil {
			panic(fmt.Errorf("failed to get minio config: %w", err))
		}

		s.minioConfig = cfg
	}

	return s.minioConfig
}

// GigaChatConfig returns the giga chat configuration.
func (s *serviceProvider) GigaChatConfig() config.GigaChatConfig {
	if s.gigaChatConfig == nil {
		cfg, err := config.NewGigaChatConfig()
		if err != nil {
			panic(fmt.Errorf("failed to get giga chat config: %w", err))
		}

		s.gigaChatConfig = cfg
	}

	return s.gigaChatConfig
}

func (s *serviceProvider) ModerationConfig() config.ModerationConfig {
	if s.moderationConfig == nil {
		cfg, err := config.NewModerationConfig()
		if err != nil {
			panic(fmt.Errorf("failed to get moderation config: %w", err))
		}
		s.moderationConfig = cfg
	}

	return s.moderationConfig
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

// RedisClient returns redis client.
func (s *serviceProvider) RedisClient(ctx context.Context) *redis.Client {
	if s.redisClient == nil {
		redisClient := redis.NewClient(&redis.Options{
			Addr:     s.RedisConfig().Addr(),
			Password: s.RedisConfig().Pass(),
			DB:       0,
		})
		if err := redisClient.Ping(ctx).Err(); err != nil {
			panic(fmt.Errorf("failed to ping state storage: %w", err))
		}

		s.redisClient = redisClient
	}

	return s.redisClient
}

// ClickHouseClient returns clickhouse client.
func (s *serviceProvider) ClickHouseClient() *gorm.DB {
	if s.clickHouseClient == nil {
		gormConfig := &gorm.Config{
			TranslateError: true,
		}
		db, err := gorm.Open(clickhouse.Open(s.ClickHouseConfig().DSN()), gormConfig)
		if err != nil {
			panic(fmt.Errorf("failed to connect to clickhouse database: %w", err))
		}

		s.clickHouseClient = db
	}

	return s.clickHouseClient
}

// S3Client returns s3 client.
func (s *serviceProvider) S3Client() s3.Client {
	if s.s3Client == nil {
		s3Client, err := minio.NewClient(
			s.MinioConfig().Endpoint(),
			s.MinioConfig().PublicEndpoint(),
			s.MinioConfig().User(),
			s.MinioConfig().Pass(),
			s.MinioConfig().UseSSL(),
			s.MinioConfig().BucketName(),
		)
		if err != nil {
			panic(fmt.Errorf("failed to connect to s3 client: %w", err))
		}

		s.s3Client = s3Client
	}

	return s.s3Client
}

// AIClient returns ai client.
func (s *serviceProvider) AIClient() ai.Client {
	if s.aiClient == nil {
		aiClient, err := gigachat.NewClient(s.GigaChatConfig().ClientSecret())
		if err != nil {
			panic(fmt.Errorf("failed to connect to AI client: %w", err))
		}
		s.aiClient = aiClient
	}

	return s.aiClient
}

// ClientRepository returns the client repository.
func (s *serviceProvider) ClientRepository() repository.ClientRepository {
	if s.clientRepository == nil {
		s.clientRepository = postgresRepository.NewClientRepository(s.DBClient())
	}

	return s.clientRepository
}

// ClientUsecase returns the client usecase.
func (s *serviceProvider) ClientUsecase() usecase.ClientUsecase {
	if s.clientUsecase == nil {
		s.clientUsecase = clientUsecase.NewUsecase(s.ClientRepository())
	}

	return s.clientUsecase
}

// ClientHandler returns the client handler.
func (s *serviceProvider) ClientHandler() *client.Handler {
	if s.clientHandler == nil {
		s.clientHandler = client.NewHandler(s.ClientUsecase())
	}

	return s.clientHandler
}

// AdvertiserRepository returns the advertisers repository.
func (s *serviceProvider) AdvertiserRepository() repository.AdvertiserRepository {
	if s.advertiserRepository == nil {
		s.advertiserRepository = postgresRepository.NewAdvertiserRepository(s.DBClient())
	}

	return s.advertiserRepository
}

// AdvertiserUsecase returns the advertisers usecase.
func (s *serviceProvider) AdvertiserUsecase() usecase.AdvertiserUsecase {
	if s.advertiserUsecase == nil {
		s.advertiserUsecase = advertiserUseCase.NewUsecase(s.AdvertiserRepository())
	}

	return s.advertiserUsecase
}

// AdvertiserHandler returns the advertisers handler.
func (s *serviceProvider) AdvertiserHandler() *advertiser.Handler {
	if s.advertiserHandler == nil {
		s.advertiserHandler = advertiser.NewHandler(s.AdvertiserUsecase())
	}

	return s.advertiserHandler
}

// MLScoreRepository returns the ml-score repository.
func (s *serviceProvider) MLScoreRepository() repository.MLScoreRepository {
	if s.mlScoreRepository == nil {
		s.mlScoreRepository = postgresRepository.NewMLScoreRepository(s.DBClient())
	}

	return s.mlScoreRepository
}

// MLScoreUsecase returns the ml-score usecase.
func (s *serviceProvider) MLScoreUsecase() usecase.MLScoreUsecase {
	if s.mlScoreUsecase == nil {
		s.mlScoreUsecase = mlScoreUsecase.NewUsecase(s.MLScoreRepository())
	}

	return s.mlScoreUsecase
}

// MLScoreHandler returns the ml-score handler.
func (s *serviceProvider) MLScoreHandler() *mlscore.Handler {
	if s.mlScoreHandler == nil {
		s.mlScoreHandler = mlscore.NewHandler(s.MLScoreUsecase())
	}

	return s.mlScoreHandler
}

// CampaignRepository returns the campaign repository.
func (s *serviceProvider) CampaignRepository() repository.CampaignRepository {
	if s.campaignRepository == nil {
		s.campaignRepository = postgresRepository.NewCampaignRepository(s.DBClient(), s.ModerationConfig())
	}

	return s.campaignRepository
}

// AdRepository returns the ad repository.
func (s *serviceProvider) AdRepository(ctx context.Context) repository.AdRepository {
	if s.adRepository == nil {
		s.adRepository = redisRepository.NewAdRepository(s.RedisClient(ctx))
	}

	return s.adRepository
}

// CampaignUsecase returns the campaign usecase.
func (s *serviceProvider) CampaignUsecase(ctx context.Context) usecase.CampaignUsecase {
	if s.campaignUsecase == nil {
		s.campaignUsecase = campaignUsecase.NewUsecase(
			s.CampaignRepository(),
			s.TimeRepository(ctx),
			s.ClientRepository(),
			s.ImpressionRepository(),
			s.ClickRepository(),
			s.MLScoreRepository(),
			s.AdvertiserRepository(),
			s.StatisticRepository(),
			s.AdRepository(ctx),
			s.S3Client(),
			s.AIClient(),
		)
	}

	return s.campaignUsecase
}

// CampaignHandler returns the campaign handler.
func (s *serviceProvider) CampaignHandler(ctx context.Context) *campaign.Handler {
	if s.campaignHandler == nil {
		s.campaignHandler = campaign.NewHandler(s.CampaignUsecase(ctx), s.TimeRepository(ctx))
	}

	return s.campaignHandler
}

// TimeRepository returns the time repository.
func (s *serviceProvider) TimeRepository(ctx context.Context) repository.TimeRepository {
	if s.timeRepository == nil {
		s.timeRepository = redisRepository.NewTimeRepository(s.RedisClient(ctx))
	}

	return s.timeRepository
}

// TimeHandler returns the time handler.
func (s *serviceProvider) TimeHandler(ctx context.Context) *time.Handler {
	if s.timeHandler == nil {
		s.timeHandler = time.NewHandler(s.TimeRepository(ctx))
	}

	return s.timeHandler
}

// ImpressionRepository returns the display repository.
func (s *serviceProvider) ImpressionRepository() repository.ImpressionRepository {
	if s.displayRepository == nil {
		s.displayRepository = clickHouseRepository.NewImpressionRepository(s.ClickHouseClient())
	}

	return s.displayRepository
}

// ClickRepository returns the click repository.
func (s *serviceProvider) ClickRepository() repository.ClickRepository {
	if s.clickRepository == nil {
		s.clickRepository = clickHouseRepository.NewClickRepository(s.ClickHouseClient())
	}

	return s.clickRepository
}

// StatisticRepository returns the statistic repository.
func (s *serviceProvider) StatisticRepository() repository.StatisticRepository {
	if s.statisticRepository == nil {
		s.statisticRepository = clickHouseRepository.NewStatisticRepository(s.ClickHouseClient())
	}

	return s.statisticRepository
}

// StatisticUsecase returns the statistic usecase.
func (s *serviceProvider) StatisticUsecase() usecase.StatisticUsecase {
	if s.statisticUsecase == nil {
		s.statisticUsecase = statisticUsecase.NewUsecase(s.StatisticRepository(), s.AdvertiserRepository(), s.CampaignRepository())
	}

	return s.statisticUsecase
}

// StatisticHandler returns the statistic handler.
func (s *serviceProvider) StatisticHandler() *statistic.Handler {
	if s.statisticHandler == nil {
		s.statisticHandler = statistic.NewHandler(s.StatisticUsecase())
	}

	return s.statisticHandler
}
