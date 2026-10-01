package dto

import "github.com/vasilcov77/user-auth/internal/domain"

type UpdateUserInput struct {
	ID    string  `json:"id"`
	Email *string `json:"email"`
	Phone *string `json:"phone"`

	IdempotencyKey string `json:"idempotency_key"`
}

func (u UpdateUserInput) Validate() error {
	if u.Email == nil && u.Phone == nil {
		return domain.ErrAllFieldsForUpdateEmpty
	}

	if u.IdempotencyKey == "" {
		return domain.ErrIdempotencyKeyRequired
	}

	return nil
}
