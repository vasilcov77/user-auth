package transaction

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vasilcov77/user-auth/pkg/postgres"
)

var (
	pool       *pgxpool.Pool
	IsUnitTest bool
)

type ctxKey struct{}

func Init(p *postgres.Pool) {
	pool = p.Pool
}

// Unit of Work
type Transaction struct {
	pgx.Tx
}

type Executor interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func TryExtractTX(ctx context.Context) Executor { //nolint:ireturn
	tx, ok := ctx.Value(ctxKey{}).(*Transaction)
	if !ok {
		return pool
	}

	return tx
}
