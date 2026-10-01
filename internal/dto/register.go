package dto

import (
	"github.com/google/uuid"
)

type RegisterOutput struct {
	ID uuid.UUID `json:"id"`
}

type RegisterInput struct {
	Email    string `json:"email" required:"true"`
	Phone    string `json:"phone"`
	Password string `json:"password" required:"true"`
	Company  string `json:"company"`
}
