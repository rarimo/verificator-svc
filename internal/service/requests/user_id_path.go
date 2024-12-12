package requests

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	val "github.com/go-ozzo/ozzo-validation/v4"
)

func GetPathUserID(r *http.Request) (userID string, err error) {
	userID = strings.ToLower(chi.URLParam(r, "user_id"))

	err = val.Errors{
		"user_id": val.Validate(userID, val.Required),
	}.Filter()
	return
}
