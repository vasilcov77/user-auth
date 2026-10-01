package postgres

import (
	"context"
	"fmt"

	"github.com/vasilcov77/user-auth/internal/domain"
	"github.com/vasilcov77/user-auth/pkg/transaction"
)

func (p *Postgres) CreateCompany(ctx context.Context, company domain.Company) error {
	const sql = `INSERT INTO company (id, name, active)
                    VALUES ($1, $2, $3)`

	args := []any{
		company.ID,
		company.Name,
		company.Active,
	}

	txOrPool := transaction.TryExtractTX(ctx)

	_, err := txOrPool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("txOrPool.Exec: %w", err)
	}

	return nil
}
