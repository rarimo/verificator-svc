package requests

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	val "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/verificator-svc/internal/service/ctx"
)

func GetPathUserID(r *http.Request) (userID string, err error) {
	userID = chi.URLParam(r, "user_id")

	if ctx.Verifiers(r).LowerCaseUserID {
		userID = strings.ToLower(userID)
	}

	err = val.Errors{
		"user_id": val.Validate(userID, val.Required),
	}.Filter()
	return
}
