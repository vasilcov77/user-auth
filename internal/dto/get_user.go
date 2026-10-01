package dto

import (
	"github.com/vasilcov77/user-auth/internal/domain"
)

type GetUserOutput struct {
	domain.User
}

type GetUserInput struct {
	ID string
}
