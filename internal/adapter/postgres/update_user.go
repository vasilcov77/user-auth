package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/vasilcov77/user-auth/internal/domain"
	"github.com/vasilcov77/user-auth/pkg/transaction"
)

func (p *Postgres) UpdateUser(ctx context.Context, user domain.User) error {
	const sql = `UPDATE "user" SET email = $1, phone = $2, verified_email = $3, verified_phone = $4, updated_at = NOW()
                     WHERE id = $5`

	args := []any{
		user.Email,
		user.Phone,
		user.VerifiedEmail,
		user.VerifiedPhone,
		user.ID,
	}

	txOrPool := transaction.TryExtractTX(ctx)

	_, err := txOrPool.Exec(ctx, sql, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrNotFound
		}

		return fmt.Errorf("txOrPool.Exec: %w", err)
	}

	return nil
}
