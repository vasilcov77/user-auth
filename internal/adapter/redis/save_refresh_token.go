package redis

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func (r *Redis) SaveRefreshToken(ctx context.Context, token string, userID uuid.UUID, ttl time.Duration) error {
	key := "refresh_token:" + token

	err := r.redis.Set(ctx, key, userID.String(), ttl).Err()
	if err != nil {
		log.Error().Err(err).Msg("redis: SaveRefreshToken: r.redis.Set")
		return err
	}

	return nil
}
