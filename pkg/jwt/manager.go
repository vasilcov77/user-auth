package jwt

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Config struct {
	Secret     string `envconfig:"JWT_ACCESS_SECRET"     required:"true"`
	AccessTTL  string `envconfig:"JWT_ACCESS_TTL"     required:"true"`
	RefreshTTL string `envconfig:"JWT_REFRESH_TTL"     required:"true"`
}

type Manager struct {
	config Config
}

func New(c Config) *Manager {
	return &Manager{config: c}
}

func (m *Manager) GenerateAccessToken(userID uuid.UUID) (string, error) {
	ttl, err := strconv.Atoi(m.config.AccessTTL)
	if err != nil {
		return "", err
	}

	now := time.Now()
	expiresAt := now.Add(time.Duration(ttl) * time.Minute)

	claims := jwt.MapClaims{
		"sub": userID.String(),
		"exp": expiresAt.Unix(),
		"iat": now.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	res, err := token.SignedString([]byte(m.config.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return res, nil
}

func (m *Manager) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (m *Manager) HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (m *Manager) ExpireRefreshToken() (time.Duration, error) {
	ttl, err := strconv.Atoi(m.config.RefreshTTL)
	if err != nil {
		return 0, fmt.Errorf("invalid refresh ttl: %w", err)
	}

	return time.Duration(ttl) * 24 * time.Hour, nil
}

func (m *Manager) AccessSecret() []byte {
	return []byte(m.config.Secret)
}
