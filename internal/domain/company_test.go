package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewCompany(t *testing.T) {
	t.Run("создание компании с валидным именем", func(t *testing.T) {
		name := "Test Company"
		company := NewCompany(name)

		require.NotEqual(t, uuid.Nil, company.ID, "ожидался уникальный ID")
		require.Equal(t, name, company.Name, "имя компании должно совпадать")
		require.True(t, company.Active, "компания должна быть активной")
	})

	t.Run("создание компании с пустым именем", func(t *testing.T) {
		name := ""
		company := NewCompany(name)

		require.Equal(t, uuid.Nil, company.ID, "ожидался нулевой UUID для невалидной компании")
		require.Empty(t, company.Name, "имя должно быть пустым")
	})

	t.Run("создание компании с nil именем", func(t *testing.T) {
		var name string
		company := NewCompany(name)

		require.Equal(t, uuid.Nil, company.ID, "ожидался нулевой UUID для невалидной компании")
	})
}

func TestValidate(t *testing.T) {
	t.Run("валидация корректной компании", func(t *testing.T) {
		company := Company{
			ID:     uuid.New(),
			Name:   "Valid Company",
			Active: true,
		}

		err := company.Validate()
		require.NoError(t, err, "валидация должна пройти успешно")
	})

	t.Run("валидация компании с пустым именем", func(t *testing.T) {
		company := Company{
			ID:     uuid.New(),
			Name:   "",
			Active: true,
		}

		err := company.Validate()
		require.Error(t, err, "ожидалась ошибка валидации для пустого имени")
	})
}
