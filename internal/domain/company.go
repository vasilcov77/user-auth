package domain

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Company struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"-"`
	Name      string    `json:"name" validate:"required"`
	Active    bool      `json:"active"`
}

func NewCompany(name string) Company {
	c := Company{
		ID:     uuid.New(),
		Name:   name,
		Active: true,
	}

	if err := c.Validate(); err != nil {
		return Company{}
	}

	return c
}

func (c Company) Validate() error {
	var validate = validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(c)
	if err != nil {
		return fmt.Errorf("validate.Company.Struct: %w", err)
	}

	return nil
}
