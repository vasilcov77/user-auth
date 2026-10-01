package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/vasilcov77/user-auth/internal/adapter/kafka_producer"
	"github.com/vasilcov77/user-auth/internal/controller/kafka_consumer"
	"github.com/vasilcov77/user-auth/internal/controller/worker"
	"github.com/vasilcov77/user-auth/pkg/httpserver"
	"github.com/vasilcov77/user-auth/pkg/jwt"
	"github.com/vasilcov77/user-auth/pkg/logger"
	"github.com/vasilcov77/user-auth/pkg/postgres"
	"github.com/vasilcov77/user-auth/pkg/redis"
)

type App struct {
	Name    string `envconfig:"APP_NAME"    required:"true"`
	Version string `envconfig:"APP_VERSION" required:"true"`
}

type Config struct {
	App           App
	HTTP          httpserver.Config
	Logger        logger.Config
	Postgres      postgres.Config
	Redis         redis.Config
	KafkaConsumer kafka_consumer.Config
	KafkaProducer kafka_producer.Config
	OutboxKafka   worker.OutboxKafkaConfig
	JWT           jwt.Config
}

func New() (Config, error) {
	var config Config

	err := godotenv.Load(".env")
	if err != nil {
		return config, fmt.Errorf("godotenv.Load: %w", err)
	}

	err = envconfig.Process("", &config)
	if err != nil {
		return config, fmt.Errorf("envconfig.Process: %w", err)
	}

	return config, nil
}
