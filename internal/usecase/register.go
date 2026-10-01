package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/vasilcov77/user-auth/internal/domain"
	"github.com/vasilcov77/user-auth/internal/dto"
	"github.com/vasilcov77/user-auth/pkg/transaction"
)

func (p *Ports) Register(ctx context.Context, input dto.RegisterInput) (dto.RegisterOutput, error) {
	var output dto.RegisterOutput

	existUser, err := p.postgres.GetUserByEmail(ctx, input.Email)
	if err != nil {
		if !errors.Is(err, domain.ErrNotFound) {
			return output, fmt.Errorf("postgres.GetUserByEmail: %w", err)
		}
	}

	if existUser.ID != uuid.Nil {
		return output, fmt.Errorf("user already exists")
	}

	user, err := domain.NewUser(input.Email, input.Phone, input.Password)
	if err != nil {
		return output, fmt.Errorf("domain.NewUser: %w", err)
	}

	event, err := user.ToEvent("user-auth-topic")
	if err != nil {
		return output, fmt.Errorf("profile.ToEvent: %w", err)
	}

	company := domain.NewCompany(input.Company)

	err = transaction.Wrap(ctx, func(ctx context.Context) error {
		err = p.postgres.CreateUser(ctx, user)
		if err != nil {
			return fmt.Errorf("postgres.CreateProfile: %w", err)
		}

		err = p.postgres.CreateCompany(ctx, company)
		if err != nil {
			return fmt.Errorf("postgres.CreateProperty: %w", err)
		}

		err = p.postgres.CreateUserCompany(ctx, user, company)
		if err != nil {
			return fmt.Errorf("postgres.CreateProperty: %w", err)
		}

		err = p.postgres.SaveOutboxKafka(ctx, event)
		if err != nil {
			return fmt.Errorf("postgres.SaveOutboxKafka: %w", err)
		}

		return nil
	})
	if err != nil {
		return output, fmt.Errorf("transaction.Wrap: %w", err)
	}

	return dto.RegisterOutput{
		ID: user.ID,
	}, nil
}
