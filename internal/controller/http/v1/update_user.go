package v1

import (
	"encoding/json"
	"net/http"

	"github.com/vasilcov77/user-auth/internal/dto"
	"github.com/vasilcov77/user-auth/pkg/render"
)

func (h *Handlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	input := dto.UpdateUserInput{}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		render.Error(w, err, http.StatusBadRequest, "json decode error")

		return
	}

	err = h.ports.UpdateUser(ctx, input)
	if err != nil {
		render.Error(w, err, http.StatusBadRequest, "request failed")

		return
	}

	w.WriteHeader(http.StatusOK)
}
