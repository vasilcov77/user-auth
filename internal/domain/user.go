package domain

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Email string

type User struct {
	ID            uuid.UUID `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     time.Time `json:"-"`
	Email         string    `json:"email" validate:"email"`
	Phone         string    `json:"phone" validate:"e164"`
	Password      string    `json:"password"`
	Active        bool      `json:"active"`
	VerifiedPhone bool      `json:"verified_phone"`
	VerifiedEmail bool      `json:"verified_email"`
}

func NewUser(email, phone, password string) (User, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("bcrypt.GenerateFromPassword: %w", err)
	}

	u := User{
		ID:            uuid.New(),
		Email:         email,
		Phone:         phone,
		Password:      string(hashedPassword),
		Active:        true,
		VerifiedPhone: false,
		VerifiedEmail: false,
	}

	if err := u.Validate(); err != nil {
		return User{}, fmt.Errorf("u.Validate: %w", err)
	}

	return u, nil
}

func (u User) Validate() error {
	var validate = validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(u)
	if err != nil {
		return fmt.Errorf("validate.User.Struct: %w", err)
	}

	return nil
}

func (u User) IsDeleted() bool {
	return !u.DeletedAt.IsZero()
}

func (u User) ToEvent(topic string) (Event, error) {
	value, err := json.Marshal(u)
	if err != nil {
		return Event{}, fmt.Errorf("json.Marshal: %w", err)
	}

	return Event{
		Topic: topic,
		Key:   []byte(u.ID.String()),
		Value: value,
	}, nil
}

func (u User) CheckPassword(plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainPassword))
}
