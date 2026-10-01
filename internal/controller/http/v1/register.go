package v1

import (
	"encoding/json"
	"net/http"

	"github.com/vasilcov77/user-auth/internal/dto"
	"github.com/vasilcov77/user-auth/pkg/render"
)

func (h *Handlers) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	input := dto.RegisterInput{}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		render.Error(w, err, http.StatusBadRequest, "json decode error")

		return
	}

	output, err := h.ports.Register(ctx, input)
	if err != nil {
		render.Error(w, err, http.StatusBadRequest, "request failed")

		return
	}

	render.JSON(w, output, http.StatusOK)
}
