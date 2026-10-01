package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/vasilcov77/user-auth/internal/controller/http/v1"
	"github.com/vasilcov77/user-auth/internal/usecase"
	"github.com/vasilcov77/user-auth/pkg/jwt"
	"github.com/vasilcov77/user-auth/pkg/logger"
)

func Routing(r *chi.Mux, ucp *usecase.Ports, mgr *jwt.Manager) {
	handlerV1 := v1.New(ucp)

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(logger.Middleware)

		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", handlerV1.Register)
			r.Post("/login", handlerV1.Login)
			r.Post("/logout", handlerV1.Logout)
			r.Post("/refresh", handlerV1.Refresh)
		})

		r.Route("/user", func(r chi.Router) {
			r.Use(jwt.Middleware(mgr))
			r.Get("/{id}", handlerV1.GetUser)
			r.Put("/", handlerV1.UpdateUser)
		})
	})
}
