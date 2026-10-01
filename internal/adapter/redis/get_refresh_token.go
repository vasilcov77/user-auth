package redis

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

func (r *Redis) GetUserIDByRefreshToken(ctx context.Context, token string) (string, error) {
	key := "refresh_token:" + token

	val, err := r.redis.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// Токен не найден или истёк
			return "", errors.New("token not found or expired")
		}
		log.Error().Err(err).Msg("redis: GetUserIDByRefreshToken: r.redis.Get")
		return "", err
	}

	return val, nil
}
