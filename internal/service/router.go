package service

import (
	"github.com/go-chi/chi"
	"github.com/rarimo/verificator-svc/internal/config"
	"github.com/rarimo/verificator-svc/internal/data/pg"
	"github.com/rarimo/verificator-svc/internal/service/ctx"
	"github.com/rarimo/verificator-svc/internal/service/handlers"
	"github.com/rarimo/verificator-svc/internal/service/middlewares"
	"gitlab.com/distributed_lab/ape"
)

func (s *service) router(cfg config.Config) chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			ctx.CtxLog(cfg.Log()),
			ctx.CtxVerifyUsersQ(pg.NewVerifyUsersQ(cfg.DB().Clone())),
			ctx.CtxVerifiers(cfg.Verifiers()),
			ctx.CtxCallback(cfg.CallbackConfig()),
			ctx.CtxSignatureVerification(cfg.SignatureVerificationConfig()),
			ctx.CtxAuthClient(cfg.Auth()),
		),
	)
	authMW := middlewares.Auth(cfg.Auth(), cfg.Log())
	r.Route("/integrations/verificator-svc", func(r chi.Router) {
		r.Route("/private", func(r chi.Router) {
			r.With(authMW).Get("/proof-parameters", handlers.GetProofParameters)
			r.Get("/proof/{user_id}", handlers.GetProofByUserID)
			r.Get("/verification-status/{user_id}", handlers.GetVerificationStatusById)
			r.With(authMW).Delete("/user/{user_id}", handlers.DeleteUser)
			r.With(authMW).Post("/verification-link", handlers.VerificationLink)
		})
		r.Route("/public", func(r chi.Router) {
			r.Post("/callback/{user_id}", handlers.VerificationCallback)
			r.Get("/proof-params/{user_id_hash}", handlers.GetProofParamsById)
		})

		r.Route("/light", func(r chi.Router) {
			r.Route("/public", func(r chi.Router) {
				r.Post("/callback-sign/{user_id_hash}", handlers.VerificationSignatureCallback)
				r.Get("/proof-params/{user_id_hash}", handlers.GetProofParamsLightById)
			})
			r.Route("/private", func(r chi.Router) {
				r.With(authMW).Post("/verification-link", handlers.VerificationLinkLight)
				r.With(authMW).Delete("/user/{user_id}", handlers.DeleteUser)
				r.Get("/user/{user_id}", handlers.GetUser)
				r.Get("/verification-status/{user_id}", handlers.GetVerificationStatusById)
			})
		})
		r.Route("/v2", func(r chi.Router) {
			r.Route("/private", func(r chi.Router) {
				r.With(authMW).Post("/verification-link", handlers.VerificationLinkV2)
			})
		})

	})

	return r
}
