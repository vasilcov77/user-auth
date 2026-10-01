package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/vasilcov77/user-auth/config"
	"github.com/vasilcov77/user-auth/internal/adapter/kafka_producer"
	"github.com/vasilcov77/user-auth/internal/adapter/postgres"
	"github.com/vasilcov77/user-auth/internal/adapter/redis"
	"github.com/vasilcov77/user-auth/internal/controller/http"
	"github.com/vasilcov77/user-auth/internal/controller/kafka_consumer"
	"github.com/vasilcov77/user-auth/internal/controller/worker"
	"github.com/vasilcov77/user-auth/internal/usecase"
	"github.com/vasilcov77/user-auth/pkg/httpserver"
	"github.com/vasilcov77/user-auth/pkg/jwt"
	pgpool "github.com/vasilcov77/user-auth/pkg/postgres"
	redislib "github.com/vasilcov77/user-auth/pkg/redis"
	"github.com/vasilcov77/user-auth/pkg/router"
	"github.com/vasilcov77/user-auth/pkg/transaction"
)

func Run(ctx context.Context, c config.Config) error {
	//Postgres
	pgPool, err := pgpool.New(ctx, c.Postgres)
	if err != nil {
		return fmt.Errorf("postgres.New: %w", err)
	}

	transaction.Init(pgPool)

	//Redis
	redisClient, err := redislib.New(c.Redis)
	if err != nil {
		return fmt.Errorf("redislib.New: %w", err)
	}

	//jwt
	jwtManager := jwt.New(c.JWT)

	// Kafka producer
	kafkaProducer := kafka_producer.New(c.KafkaProducer)

	// UseCasePorts
	ucp := usecase.New(
		postgres.New(),
		redis.New(redisClient),
		kafkaProducer,
		jwtManager,
	)

	// Kafka consumer
	kafkaConsumer := kafka_consumer.New(c.KafkaConsumer, ucp)

	// Outbox Kafka worker
	outboxKafkaWorker := worker.NewOutboxKafka(ucp, c.OutboxKafka)

	// HTTP
	r := router.New()
	http.Routing(r, ucp, jwtManager)
	httpServer := httpserver.New(r, c.HTTP)

	log.Info().Msg("app: started")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig // wait signal

	log.Info().Msg("app: got signal to stop")

	// Controllers close
	httpServer.Close()
	outboxKafkaWorker.Close()
	kafkaConsumer.Close()

	// Clients close
	redisClient.Close()
	pgPool.Close()

	log.Info().Msg("app: stopped")

	return nil
}
