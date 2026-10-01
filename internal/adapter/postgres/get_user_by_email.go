package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/vasilcov77/user-auth/internal/domain"
	"github.com/vasilcov77/user-auth/pkg/transaction"
)

func (p *Postgres) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	const sql = `SELECT ID, email, password, active
                    FROM "user" WHERE email = $1`

	dto := struct {
		ID       pgtype.Text
		Email    pgtype.Text
		Password pgtype.Text
		Active   pgtype.Bool
	}{}

	dest := []any{
		&dto.ID,
		&dto.Email,
		&dto.Password,
		&dto.Active,
	}

	txOrPool := transaction.TryExtractTX(ctx)

	err := txOrPool.QueryRow(ctx, sql, email).Scan(dest...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrNotFound
		}

		return domain.User{}, fmt.Errorf("txOrPool.QueryRow: %w", err)
	}

	id, err := uuid.Parse(dto.ID.String)
	if err != nil {
		return domain.User{}, domain.ErrUUIDInvalid
	}

	user := domain.User{
		ID:       id,
		Email:    dto.Email.String,
		Password: dto.Password.String,
		Active:   dto.Active.Bool,
	}

	return user, nil
}
