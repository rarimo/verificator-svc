package middlewares

import (
	"net/http"

	"github.com/rarimo/verificator-svc/internal/service/handlers"
	"github.com/rarimo/web3-auth-svc/pkg/auth"
	"github.com/rarimo/web3-auth-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
)

func Auth(auth *auth.Client, log *logan.Entry) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !auth.Enabled {
				log.Warn("authentication is disabled, validation skipped")
				next.ServeHTTP(w, r.WithContext(handlers.CtxUserClaims([]resources.Claim{})(r.Context())))
				return
			}

			claims, err := auth.ValidateJWT(r)
			if err != nil {
				log.WithError(err).Error("Got invalid auth or validation error")
				ape.RenderErr(w, problems.Unauthorized())
				return
			}

			if len(claims) == 0 {
				log.Debug("Claims are empty")
				ape.RenderErr(w, problems.Unauthorized())
				return
			}

			next.ServeHTTP(w, r.WithContext(handlers.CtxUserClaims(claims)(r.Context())))
		})
	}
}
