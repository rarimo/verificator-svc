package middlewares

import (
	"net/http"

	"github.com/rarimo/web3-auth-svc/pkg/auth"
	"github.com/rarimo/web3-auth-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type AuthOption func(claims ...resources.Claim) error
type AuthMiddleware func(opts ...AuthOption) func(http.Handler) http.Handler

var ErrUserIsNotAdmin = errors.New("user is not an admin")

func Auth(auth *auth.Client, log *logan.Entry) AuthMiddleware {
	return func(opts ...AuthOption) func(http.Handler) http.Handler {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				claims, err := auth.ValidateJWT(r)
				if err != nil {
					log.WithError(err).Error("Got invalid auth or validation error")
					ape.RenderErr(w, problems.Unauthorized())
					return
				}

				if len(claims) == 0 {
					log.Error("Claims are empty")
					ape.RenderErr(w, problems.Unauthorized())
					return
				}

				for _, option := range opts {
					if err = option(claims...); err != nil {
						log.WithError(err).Error("Failed to validate option")
						ape.RenderErr(w, problems.Unauthorized())
						return
					}
				}

				next.ServeHTTP(w, r.WithContext(r.Context()))
			})
		}
	}
}

func AdminGrant(claims ...resources.Claim) error {
	if !auth.Authenticates(claims, auth.AdminGrant) {
		return ErrUserIsNotAdmin
	}

	return nil
}
