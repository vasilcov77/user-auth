package usecase

import (
	"context"
	"fmt"

	"github.com/vasilcov77/user-auth/internal/dto"
)

func (p *Ports) Logout(ctx context.Context, input dto.LogoutInput) error {
	_, err := p.redis.GetUserIDByRefreshToken(ctx, p.jwt.HashToken(input.RefreshToken))
	if err != nil {
		return fmt.Errorf("redis.GetUserIDByRefreshToken: %w", err)
	}

	err = p.redis.DeleteRefreshToken(ctx, p.jwt.HashToken(input.RefreshToken))
	if err != nil {
		return fmt.Errorf("redis.DeleteRefreshToken: %w", err)
	}

	return nil
}
