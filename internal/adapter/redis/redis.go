package redis

import (
	"time"

	"github.com/vasilcov77/user-auth/pkg/redis"
)

const (
	idempotencyPrefix = "idempotency:"
	ttl               = time.Hour
)

type Redis struct {
	redis *redis.Client
}

func New(client *redis.Client) *Redis {
	return &Redis{
		redis: client,
	}
}
