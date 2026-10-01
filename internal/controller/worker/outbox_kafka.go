package worker

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/vasilcov77/user-auth/internal/usecase"
)

type OutboxKafkaConfig struct {
	Limit int `default:"10" envconfig:"OUTBOX_KAFKA_WORKER_LIMIT"`
}

type OutboxKafka struct {
	config OutboxKafkaConfig
	ports  *usecase.Ports
	stop   chan struct{}
	done   chan struct{}
}

func NewOutboxKafka(ucp *usecase.Ports, c OutboxKafkaConfig) *OutboxKafka {
	w := &OutboxKafka{
		config: c,
		ports:  ucp,
		stop:   make(chan struct{}),
		done:   make(chan struct{}),
	}

	go w.run()

	return w
}

func (w *OutboxKafka) run() {
	log.Info().Msg("outbox kafka worker: started")

FOR:
	for {
		count, err := w.ports.OutboxReadAndProduce(context.Background(), w.config.Limit)
		if err != nil {
			log.Error().Err(err).Msg("outbox kafka worker: read and produce failed")
		}

		log.Debug().Int("count", count).Msg("outbox kafka worker: read and produce")

		var duration time.Duration

		if count < w.config.Limit {
			duration = 10 * time.Second

			log.Debug().Msg("outbox kafka worker: sleeping 10s")
		}

		select {
		case <-w.stop:
			break FOR
		case <-time.After(duration):
		}
	}

	close(w.done)
}

func (w *OutboxKafka) Close() {
	log.Info().Msg("outbox kafka worker: closing")

	close(w.stop)

	<-w.done

	log.Info().Msg("outbox kafka worker: closed")
}
