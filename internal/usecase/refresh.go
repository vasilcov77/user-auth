package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/vasilcov77/user-auth/internal/dto"
)

func (p *Ports) Refresh(ctx context.Context, input dto.RefreshInput) (dto.RefreshOutput, error) {
	var output dto.RefreshOutput

	userIdString, err := p.redis.GetUserIDByRefreshToken(ctx, p.jwt.HashToken(input.RefreshToken))
	if err != nil {
		return output, fmt.Errorf("redis.GetUserIDByRefreshToken: %w", err)
	}

	userId, err := uuid.Parse(userIdString)
	if err != nil {
		return output, fmt.Errorf("uuid.Parse: %w", err)
	}

	accessToken, err := p.jwt.GenerateAccessToken(userId)
	if err != nil {
		return output, fmt.Errorf("jwt.GenerateAccessToken: %w", err)
	}

	refreshToken, err := p.jwt.GenerateRefreshToken()
	if err != nil {
		return output, fmt.Errorf("jwt.GenerateAccessToken: %w", err)
	}

	expireToken, err := p.jwt.ExpireRefreshToken()
	if err != nil {
		return output, fmt.Errorf("jwt.ExpireRefreshToken: %w", err)
	}

	err = p.redis.SaveRefreshToken(ctx, p.jwt.HashToken(refreshToken), userId, expireToken)
	if err != nil {
		return output, fmt.Errorf("redis.SaveRefreshToken: %w", err)
	}

	err = p.redis.DeleteRefreshToken(ctx, p.jwt.HashToken(input.RefreshToken))
	if err != nil {
		return output, fmt.Errorf("redis.DeleteRefreshToken: %w", err)
	}

	return dto.RefreshOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
