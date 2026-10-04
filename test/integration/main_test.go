//go:build integration

package test

import (
	"context"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/vasilcov77/user-auth/config"
	"github.com/vasilcov77/user-auth/internal/adapter/kafka_producer"
	"github.com/vasilcov77/user-auth/internal/app"
	"github.com/vasilcov77/user-auth/internal/controller/kafka_consumer"
	"github.com/vasilcov77/user-auth/internal/controller/worker"
	"github.com/vasilcov77/user-auth/pkg/httpclient"
	"github.com/vasilcov77/user-auth/pkg/httpserver"
	"github.com/vasilcov77/user-auth/pkg/jwt"
	"github.com/vasilcov77/user-auth/pkg/postgres"
	"github.com/vasilcov77/user-auth/pkg/redis"
)

// Prepare:  make up
// Run test: make integration-test

var ctx = context.Background()

func Test_Integration(t *testing.T) {
	suite.Run(t, &Suite{})
}

type Suite struct {
	suite.Suite
	*require.Assertions

	client *httpclient.Client
}

func (s *Suite) SetupSuite() { // В начале всех тестов
	s.Assertions = s.Require()

	s.ResetMigrations()

	// Config
	c := config.Config{
		App: config.App{
			Name:    "user-auth",
			Version: "test",
		},
		HTTP: httpserver.Config{
			Port: "8080",
		},
		Postgres: postgres.Config{
			Host:     "localhost",
			Port:     "5432",
			User:     "login",
			Password: "pass",
			DBName:   "postgres",
		},
		Redis: redis.Config{
			Addr: "localhost:6379",
		},
		KafkaConsumer: kafka_consumer.Config{
			Addr:  []string{"localhost:9092"},
			Topic: "user-auth-topic",
			Group: "user-auth-group",
		},
		KafkaProducer: kafka_producer.Config{
			Addr: []string{"localhost:9092"},
		},
		OutboxKafka: worker.OutboxKafkaConfig{
			Limit: 10,
		},
		JWT: jwt.Config{
			Secret:     "SuPSEyNgIvm9oANzmkwWQK8znJoQfQp5rgubNfF72e4",
			AccessTTL:  "15",
			RefreshTTL: "7",
		},
	}

	// Logger and OTEL disable
	log.Logger = zerolog.Nop()

	// Server
	go func() {
		err := app.Run(context.Background(), c)
		s.NoError(err)
	}()

	// API client
	s.client = httpclient.New(httpclient.Config{Host: "localhost", Port: "8080"})

	time.Sleep(1 * time.Second)
}

func (s *Suite) TearDownSuite() {} // В конце всех тестов

func (s *Suite) SetupTest() {} // Перед каждым тестом

func (s *Suite) TearDownTest() {} // После каждого теста
