package postgres

import (
	"context"
	"fmt"

	"github.com/vasilcov77/user-auth/internal/domain"
	"github.com/vasilcov77/user-auth/pkg/transaction"
)

func (p *Postgres) CreateUserCompany(ctx context.Context, user domain.User, company domain.Company) error {
	const sql = `INSERT INTO user_company (user_id, company_id)
                    VALUES ($1, $2)`

	args := []any{
		user.ID,
		company.ID,
	}

	txOrPool := transaction.TryExtractTX(ctx)

	_, err := txOrPool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("txOrPool.Exec: %w", err)
	}

	return nil
}
