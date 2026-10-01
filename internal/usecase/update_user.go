package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/vasilcov77/user-auth/internal/domain"
	"github.com/vasilcov77/user-auth/internal/dto"
	"github.com/vasilcov77/user-auth/pkg/transaction"
)

func (p *Ports) UpdateUser(ctx context.Context, input dto.UpdateUserInput) error {
	// Если мы уже обработали этот запрос, просто выходим
	if p.redis.IsExists(ctx, input.IdempotencyKey) {
		log.Info().Msg("usecase UpdateUser: idempotent request - skipping")

		return nil
	}

	err := input.Validate()
	if err != nil {
		return fmt.Errorf("input.Validate: %w", err)
	}

	id, err := uuid.Parse(input.ID)
	if err != nil {
		return domain.ErrUUIDInvalid
	}

	err = transaction.Wrap(ctx, func(ctx context.Context) error {
		user, err := p.postgres.GetUser(ctx, id)
		if err != nil {
			return fmt.Errorf("postgres.GetUser: %w", err)
		}

		if user.IsDeleted() {
			return domain.ErrNotFound
		}

		newUser := update(user, input)

		if newUser == user {
			return nil
		}

		err = p.postgres.UpdateUser(ctx, newUser)
		if err != nil {
			return fmt.Errorf("postgres.UpdateUser: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("transaction.Wrap: %w", err)
	}

	return nil
}

func update(user domain.User, input dto.UpdateUserInput) domain.User {
	if input.Email != nil {
		user.Email = *input.Email
		user.VerifiedEmail = false
	}

	if input.Phone != nil {
		user.Phone = *input.Phone
		user.VerifiedPhone = false
	}

	return user
}
