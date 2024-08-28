package service

import (
	"github.com/go-chi/chi"
	"github.com/rarimo/verificator-svc/internal/config"
	"github.com/rarimo/verificator-svc/internal/data/pg"
	"github.com/rarimo/verificator-svc/internal/service/handlers"
	"gitlab.com/distributed_lab/ape"
)

func (s *service) router(cfg config.Config) chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			handlers.CtxLog(cfg.Log()),
			handlers.CtxVerifyUsersQ(pg.NewVerifyUsersQ(cfg.DB().Clone())),
			handlers.CtxVerifiers(cfg.Verifiers()),
			handlers.CtxCallback(cfg.CallbackConfig()),
			handlers.CtxProofParameters(cfg.ProofParametersConfig()),
		),
	)
	r.Route("/integrations/verificator-svc", func(r chi.Router) {
		r.Route("/private", func(r chi.Router) {
			r.Get("/proof-parameters", handlers.GetProofParameters)
			r.Get("/verification-status/{user_id}", handlers.GetVerificationStatusById)
		})
		r.Route("/public", func(r chi.Router) {
			r.Post("/callback/{user_id}", handlers.VerificationCallback)
		})
	})

	return r
}
