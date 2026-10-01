package v1

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vasilcov77/user-auth/internal/domain"
	"github.com/vasilcov77/user-auth/internal/dto"
	"github.com/vasilcov77/user-auth/pkg/render"
)

func (h *Handlers) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	input := dto.GetUserInput{
		ID: chi.URLParam(r, "id"),
	}

	output, err := h.ports.GetUser(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			render.Error(w, err, http.StatusNotFound, "request failed")

		default:
			render.Error(w, err, http.StatusBadRequest, "request failed")
		}

		return
	}

	render.JSON(w, output, http.StatusOK)
}
