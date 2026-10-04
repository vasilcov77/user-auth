// user_test.go
package domain

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestNewUser(t *testing.T) {
	t.Run("создание пользователя с валидными данными", func(t *testing.T) {
		email := "test@example.com"
		phone := "+79001234567"
		password := "securePassword123"

		user, err := NewUser(email, phone, password)

		require.NoError(t, err, "NewUser не должен возвращать ошибку для валидных данных")
		require.NotEmpty(t, user.ID, "ID пользователя не должен быть пустым")
		require.Equal(t, email, user.Email, "Email должен совпадать с переданным")
		require.Equal(t, phone, user.Phone, "Phone должен совпадать с переданным")
		require.True(t, user.Active, "Active должен быть true по умолчанию")
		require.False(t, user.VerifiedPhone, "VerifiedPhone должен быть false по умолчанию")
		require.False(t, user.VerifiedEmail, "VerifiedEmail должен быть false по умолчанию")
		require.NotEmpty(t, user.Password, "Password не должен быть пустым (должен быть хеширован)")
		require.False(t, user.IsDeleted(), "IsDeleted должен возвращать false для нового пользователя")
	})

	t.Run("создание пользователя с невалидным email", func(t *testing.T) {
		email := "invalid-email"
		phone := "+79001234567"
		password := "securePassword123"

		user, err := NewUser(email, phone, password)

		require.Error(t, err, "NewUser должен возвращать ошибку для невалидного email")
		require.Empty(t, user.ID, "ID пользователя должен быть пустым при ошибке")
	})

	t.Run("создание пользователя с невалильным телефоном", func(t *testing.T) {
		email := "test@example.com"
		phone := "invalid-phone"
		password := "securePassword123"

		user, err := NewUser(email, phone, password)

		require.Error(t, err, "NewUser должен возвращать ошибку для невалильного телефона")
		require.Empty(t, user.ID, "ID пользователя должен быть пустым при ошибке")
	})
}

func TestUserValidate(t *testing.T) {
	t.Run("валидация с валидными данными", func(t *testing.T) {
		user := User{
			ID:            uuid.New(),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			Email:         "test@example.com",
			Phone:         "+79001234567",
			Password:      "hashedPassword",
			Active:        true,
			VerifiedPhone: false,
			VerifiedEmail: false,
		}

		err := user.Validate()

		require.NoError(t, err, "Validate не должен возвращать ошибку для валидных данных")
	})

	t.Run("валидация с невалильным email", func(t *testing.T) {
		user := User{
			ID:            uuid.New(),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			Email:         "invalid-email",
			Phone:         "+79001234567",
			Password:      "hashedPassword",
			Active:        true,
			VerifiedPhone: false,
			VerifiedEmail: false,
		}

		err := user.Validate()

		require.Error(t, err, "Validate должен возвращать ошибку для невалильного email")
	})

	t.Run("валидация с пустым телефоном", func(t *testing.T) {
		user := User{
			ID:            uuid.New(),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			Email:         "test@example.com",
			Phone:         "",
			Password:      "hashedPassword",
			Active:        true,
			VerifiedPhone: false,
			VerifiedEmail: false,
		}

		err := user.Validate()

		require.Error(t, err, "Validate должен возвращать ошибку для пустого телефона")
	})
}

func TestUserIsDeleted(t *testing.T) {
	t.Run("новый пользователь не удален", func(t *testing.T) {
		user := User{
			ID:            uuid.New(),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			Email:         "test@example.com",
			Phone:         "+79001234567",
			Password:      "hashedPassword",
			Active:        true,
			VerifiedPhone: false,
			VerifiedEmail: false,
		}

		require.False(t, user.IsDeleted(), "IsDeleted должен возвращать false для нового пользователя")
	})

	t.Run("удаленный пользователь помечен как удаленный", func(t *testing.T) {
		user := User{
			ID:            uuid.New(),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			DeletedAt:     time.Now(), // Устанавливаем DeletedAt
			Email:         "test@example.com",
			Phone:         "+79001234567",
			Password:      "hashedPassword",
			Active:        true,
			VerifiedPhone: false,
			VerifiedEmail: false,
		}

		require.True(t, user.IsDeleted(), "IsDeleted должен возвращать true для удаленного пользователя")
	})
}

func TestUserCheckPassword(t *testing.T) {
	t.Run("проверка правильного пароля", func(t *testing.T) {
		password := "securePassword123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		user := User{
			ID:            uuid.New(),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			Email:         "test@example.com",
			Phone:         "+79001234567",
			Password:      string(hashedPassword),
			Active:        true,
			VerifiedPhone: false,
			VerifiedEmail: false,
		}

		err := user.CheckPassword(password)

		require.NoError(t, err, "CheckPassword не должен возвращать ошибку для правильного пароля")
	})

	t.Run("проверка неправильного пароля", func(t *testing.T) {
		password := "securePassword123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		user := User{
			ID:            uuid.New(),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			Email:         "test@example.com",
			Phone:         "+79001234567",
			Password:      string(hashedPassword),
			Active:        true,
			VerifiedPhone: false,
			VerifiedEmail: false,
		}

		err := user.CheckPassword("wrongPassword")

		require.Error(t, err, "CheckPassword должен возвращать ошибку для неправильного пароля")
	})
}

func TestUserToEvent(t *testing.T) {
	t.Run("создание события из пользователя", func(t *testing.T) {
		email := "test@example.com"
		phone := "+79001234567"
		password := "securePassword123"

		user, _ := NewUser(email, phone, password)
		topic := "user.created"

		event, err := user.ToEvent(topic)

		require.NoError(t, err, "ToEvent не должен возвращать ошибку")
		require.Equal(t, topic, event.Topic, "Topic события должен совпадать с переданным")
		require.Equal(t, []byte(user.ID.String()), event.Key, "Key события должен быть ID пользователя")
		require.NotEmpty(t, event.Value, "Value события не должно быть пустым")

		// Проверяем, что Value можно распарсить обратно в User
		var unmarshaledUser User
		err = json.Unmarshal(event.Value, &unmarshaledUser)
		require.NoError(t, err, "Value события должен быть валидным JSON")
		require.Equal(t, user.ID, unmarshaledUser.ID, "ID пользователя в событии должен совпадать с исходным")
		require.Equal(t, user.Email, unmarshaledUser.Email, "Email пользователя в событии должен совпадать с исходным")
	})
}
