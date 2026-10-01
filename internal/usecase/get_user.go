package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/vasilcov77/user-auth/internal/domain"
	"github.com/vasilcov77/user-auth/internal/dto"
)

func (p *Ports) GetUser(ctx context.Context, input dto.GetUserInput) (dto.GetUserOutput, error) {
	var output dto.GetUserOutput

	id, err := uuid.Parse(input.ID)
	if err != nil {
		return output, domain.ErrUUIDInvalid
	}

	user, err := p.postgres.GetUser(ctx, id)
	if err != nil {
		return output, fmt.Errorf("postgres.GetUser: %w", err)
	}

	if user.IsDeleted() {
		return output, domain.ErrNotFound
	}

	return dto.GetUserOutput{
		User: user,
	}, nil
}
