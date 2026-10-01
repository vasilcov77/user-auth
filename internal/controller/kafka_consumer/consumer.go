package kafka_consumer

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"github.com/vasilcov77/user-auth/internal/usecase"
	"github.com/vasilcov77/user-auth/pkg/logger"
)

type Config struct {
	Addr  []string `envconfig:"KAFKA_CONSUMER_ADDR" required:"true"`
	Topic string   `default:"user-auth-topic"         envconfig:"KAFKA_CONSUMER_TOPIC"`
	Group string   `default:"user-auth-group"         envconfig:"KAFKA_CONSUMER_GROUP"`
}

type Consumer struct {
	config Config
	reader *kafka.Reader
	ports  *usecase.Ports
	stop   context.CancelFunc
	done   chan struct{}
}

func New(cfg Config, ucp *usecase.Ports) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cfg.Addr,
		Topic:          cfg.Topic,
		GroupID:        cfg.Group,
		ErrorLogger:    logger.ErrorLogger(),
		CommitInterval: 100 * time.Millisecond,
	})

	ctx, stop := context.WithCancel(context.Background())

	c := &Consumer{
		config: cfg,
		reader: r,
		ports:  ucp,
		stop:   stop,
		done:   make(chan struct{}),
	}

	go c.run(ctx)

	return c
}

func (c *Consumer) run(ctx context.Context) {
	log.Info().Msg("kafka consumer: started")

FOR:
	for {
		// Читаем сообщение из Kafka
		m, err := c.reader.FetchMessage(ctx)
		if err != nil {
			switch {
			case errors.Is(err, context.Canceled):
				log.Info().Msg("kafka consumer: context canceled")

				break FOR
			case errors.Is(err, io.EOF):
				log.Warn().Err(err).Msg("kafka consumer: FetchMessage")

				break FOR
			}

			log.Error().Err(err).Msg("kafka consumer: FetchMessage")
		}

		log.Info().Str("key", string(m.Key)).Msg("kafka consumer: message received")

		// Тут вызываем метод из usecase для обработки сообщения

		// Коммитим оффсет в consumer group
		if err = c.reader.CommitMessages(ctx, m); err != nil {
			log.Error().Err(err).Msg("kafka consumer: CommitMessages")
		}
	}

	close(c.done)
}

func (c *Consumer) Close() {
	log.Info().Msg("kafka consumer: closing")

	c.stop()

	if err := c.reader.Close(); err != nil {
		log.Error().Err(err).Msg("kafka consumer: reader.Close")
	}

	<-c.done

	log.Info().Msg("kafka consumer: closed")
}
