package v1

import (
	"github.com/vasilcov77/user-auth/internal/usecase"
)

type Handlers struct {
	ports *usecase.Ports
}

func New(ucp *usecase.Ports) *Handlers {
	return &Handlers{
		ports: ucp,
	}
}
