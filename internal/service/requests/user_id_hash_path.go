package requests

import (
	"net/http"

	"github.com/go-chi/chi"
	val "github.com/go-ozzo/ozzo-validation/v4"
)

func GetPathUserIDHash(r *http.Request) (userIdHash string, err error) {
	userIdHash = chi.URLParam(r, "user_id_hash")

	err = val.Errors{
		"user_id": val.Validate(userIdHash, val.Required),
	}.Filter()
	return
}
