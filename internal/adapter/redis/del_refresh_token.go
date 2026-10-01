package redis

import (
	"context"

	"github.com/rs/zerolog/log"
)

func (r *Redis) DeleteRefreshToken(ctx context.Context, token string) error {
	key := "refresh_token:" + token

	err := r.redis.Del(ctx, key).Err()
	if err != nil {
		log.Error().Err(err).Msg("redis: DeleteRefreshToken: r.redis.Del")
		return err
	}

	return nil
}
