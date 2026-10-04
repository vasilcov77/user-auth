package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) Login(ctx context.Context, email string, password string) (string, string, error) {
	const uri = "api/v1/auth/login"

	path := fmt.Sprintf("http://%s/%s", c.host, uri)

	request := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    email,
		Password: password,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return "", "", fmt.Errorf("json.Marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, path, bytes.NewReader(body))
	if err != nil {
		return "", "", fmt.Errorf("http.NewRequest: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("client.Do: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", "", fmt.Errorf("failed to read response body: %w", err)
		}
		return "", "", fmt.Errorf("request failed: status: %s, body: %s", resp.Status, string(respBody))
	}

	response := struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{}

	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", "", fmt.Errorf("json.Decode: %w", err)
	}

	return response.AccessToken, response.RefreshToken, nil
}
