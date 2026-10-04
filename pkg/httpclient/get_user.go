package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

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

func (c *Client) GetUser(ctx context.Context, id uuid.UUID, token string) (User, error) {
	const uri = "api/v1/user"

	path := fmt.Sprintf("http://%s/%s/%s", c.host, uri, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, http.NoBody)
	if err != nil {
		return User{}, fmt.Errorf("http.NewRequest: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.client.Do(req)
	if err != nil {
		return User{}, fmt.Errorf("client.Do: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return User{}, fmt.Errorf("io.ReadAll: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return User{}, ErrNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return User{}, fmt.Errorf("request failed: status: %s, body:%s", resp.Status, body)
	}

	var user User

	if err = json.Unmarshal(body, &user); err != nil {
		return User{}, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return user, nil
}
