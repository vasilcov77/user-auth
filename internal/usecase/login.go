package usecase

import (
	"context"
	"fmt"

	"github.com/vasilcov77/user-auth/internal/dto"
)

func (p *Ports) Login(ctx context.Context, input dto.LoginInput) (dto.LoginOutput, error) {
	var output dto.LoginOutput

	user, err := p.postgres.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return output, fmt.Errorf("postgres.GetUserByLogin: %w", err)
	}

	err = user.CheckPassword(input.Password)
	if err != nil {
		return output, fmt.Errorf("user.CheckPassword: %w", err)
	}

	accessToken, err := p.jwt.GenerateAccessToken(user.ID)
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

	err = p.redis.SaveRefreshToken(ctx, p.jwt.HashToken(refreshToken), user.ID, expireToken)
	if err != nil {
		return output, fmt.Errorf("redis.SaveRefreshToken: %w", err)
	}

	return dto.LoginOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
