package requests

import (
	"github.com/go-chi/chi"
	val "github.com/go-ozzo/ozzo-validation/v4"
	"net/http"
)

func GetUserByID(r *http.Request) (userId string, err error) {
	userId = chi.URLParam(r, "user_id")

	err = val.Errors{
		"user_id": val.Validate(userId, val.Required),
	}.Filter()
	return
}
