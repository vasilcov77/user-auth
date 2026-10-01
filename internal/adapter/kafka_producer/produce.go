package kafka_producer

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"github.com/vasilcov77/user-auth/internal/domain"
	"github.com/vasilcov77/user-auth/pkg/logger"
)

type Config struct {
	Addr []string `envconfig:"KAFKA_WRITER_ADDR" required:"true"`
}

type Producer struct {
	config Config
	writer *kafka.Writer
}

func New(c Config) *Producer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(c.Addr...),
		RequiredAcks: kafka.RequireAll,
		ErrorLogger:  logger.ErrorLogger(),
		Async:        true,
	}

	return &Producer{
		config: c,
		writer: w,
	}
}

func (p *Producer) Produce(ctx context.Context, events ...domain.Event) error {
	var msgs []kafka.Message

	for _, e := range events {
		msg := kafka.Message{
			Topic: e.Topic,
			Key:   e.Key,
			Value: e.Value,
		}

		msgs = append(msgs, msg)
	}

	err := p.writer.WriteMessages(ctx, msgs...)
	if err != nil {

		return fmt.Errorf("p.writer.WriteMessages: %w", err)
	}

	return nil
}

func (p *Producer) Close() {
	err := p.writer.Close()
	if err != nil {
		log.Error().Err(err).Msg("kafka producer: p.writer.Close")
	}

	log.Info().Msg("kafka producer: closed")
}
