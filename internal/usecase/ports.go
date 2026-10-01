package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vasilcov77/user-auth/internal/domain"
)

//go:generate mockery

type Redis interface {
	IsExists(ctx context.Context, idempotencyKey string) bool
	SaveRefreshToken(ctx context.Context, token string, userID uuid.UUID, ttl time.Duration) error
	GetUserIDByRefreshToken(ctx context.Context, token string) (string, error)
	DeleteRefreshToken(ctx context.Context, token string) error
}

type Postgres interface {
	CreateUser(ctx context.Context, user domain.User) error
	CreateCompany(ctx context.Context, user domain.Company) error
	CreateUserCompany(ctx context.Context, user domain.User, company domain.Company) error
	GetUser(ctx context.Context, userID uuid.UUID) (domain.User, error)
	UpdateUser(ctx context.Context, user domain.User) error
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)

	ReadOutboxKafka(ctx context.Context, limit int) ([]domain.Event, error)
	SaveOutboxKafka(ctx context.Context, events ...domain.Event) error
}

type Kafka interface {
	Produce(ctx context.Context, events ...domain.Event) error
}

type JWT interface {
	GenerateAccessToken(userID uuid.UUID) (string, error)
	GenerateRefreshToken() (string, error)
	AccessSecret() []byte
	HashToken(token string) string
	ExpireRefreshToken() (time.Duration, error)
}

type Ports struct {
	postgres Postgres
	redis    Redis
	kafka    Kafka
	jwt      JWT
}

func New(postgres Postgres, redis Redis, kafka Kafka, jwt JWT) *Ports {
	return &Ports{
		postgres: postgres,
		redis:    redis,
		kafka:    kafka,
		jwt:      jwt,
	}
}
