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

func (p *Postgres) GetUser(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	const sql = `SELECT created_at, updated_at, deleted_at, email, phone, active, verified_phone, verified_email
                    FROM "user" WHERE id = $1`

	dto := struct {
		CreatedAt pgtype.Timestamptz
		UpdatedAt pgtype.Timestamptz
		DeletedAt pgtype.Timestamptz
		Email     pgtype.Text
		Phone     pgtype.Text

		Active        pgtype.Bool
		VerifiedEmail pgtype.Bool
		VerifiedPhone pgtype.Bool
	}{}

	dest := []any{
		&dto.CreatedAt,
		&dto.UpdatedAt,
		&dto.DeletedAt,
		&dto.Email,
		&dto.Phone,
		&dto.Active,
		&dto.VerifiedEmail,
		&dto.VerifiedPhone,
	}

	txOrPool := transaction.TryExtractTX(ctx)

	err := txOrPool.QueryRow(ctx, sql, userID).Scan(dest...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrNotFound
		}

		return domain.User{}, fmt.Errorf("txOrPool.QueryRow: %w", err)
	}

	user := domain.User{
		ID:            userID,
		CreatedAt:     dto.CreatedAt.Time,
		UpdatedAt:     dto.UpdatedAt.Time,
		DeletedAt:     dto.DeletedAt.Time,
		Email:         dto.Email.String,
		Password:      "",
		Phone:         dto.Phone.String,
		Active:        dto.Active.Bool,
		VerifiedEmail: dto.VerifiedEmail.Bool,
		VerifiedPhone: dto.VerifiedPhone.Bool,
	}

	return user, nil
}
