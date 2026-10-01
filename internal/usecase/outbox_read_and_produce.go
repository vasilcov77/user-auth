package usecase

import (
	"context"
	"fmt"

	"github.com/vasilcov77/user-auth/pkg/transaction"
)

func (p *Ports) OutboxReadAndProduce(ctx context.Context, limit int) (count int, err error) {
	// В транзакции
	err = transaction.Wrap(ctx, func(ctx context.Context) error {
		// Читаем сообщения из outbox таблицы БД
		msg, err := p.postgres.ReadOutboxKafka(ctx, limit)
		if err != nil {
			return fmt.Errorf("u.postgres.ReadOutboxKafka: %w", err)
		}

		count = len(msg)

		// Пишем в Kafka
		err = p.kafka.Produce(ctx, msg...)
		if err != nil {
			return fmt.Errorf("u.kafka.Produce: %w", err)
		}

		return nil
	})
	if err != nil {
		return count, fmt.Errorf("transaction.Wrap: %w", err)
	}

	return count, nil
}
