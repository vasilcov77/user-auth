package postgres

import (
	"context"
	"fmt"

	"github.com/vasilcov77/user-auth/internal/domain"
	"github.com/vasilcov77/user-auth/pkg/transaction"
)

func (p *Postgres) CreateUser(ctx context.Context, user domain.User) error {
	const sql = `INSERT INTO "user" (id, email, password, phone, active, verified_email, verified_phone)
                    VALUES ($1, $2, $3, $4, $5, $6 ,$7)`

	args := []any{
		user.ID,
		user.Email,
		user.Password,
		user.Phone,
		user.Active,
		user.VerifiedEmail,
		user.VerifiedPhone,
	}

	txOrPool := transaction.TryExtractTX(ctx)

	_, err := txOrPool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("txOrPool.Exec: %w", err)
	}

	return nil
}
