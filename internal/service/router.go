package service

import (
	"github.com/go-chi/chi"
	"github.com/rarimo/verificator-svc/internal/service/handlers"
	"gitlab.com/distributed_lab/ape"
)

func (s *service) router() chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			handlers.CtxLog(s.log),
			//handlers.VerifyUsersQ()
		),
	)
	r.Route("/integrations/verificator-svc", func(r chi.Router) {
		r.Route("/public", func(r chi.Router) {
			r.Get("/proof-parameters", handlers.GetProofParameters)
			r.Get("/verification-status/{user_id}", handlers.GetVerificationStatusById)
			r.Post("/callback/{user_id}", handlers.VerificationCallback)
		})
	})

	return r
}
