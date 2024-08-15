package requests

import (
	"net/http"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func GetVerificationStatusByID(r *http.Request) (userId string, err error) {
	userId = chi.URLParam(r, "user_id")

	err = validation.Errors{
		"user_id": validation.Validate(userId, validation.Required, is.Email),
	}.Filter()
	return
}
