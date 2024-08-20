package requests

import (
	"fmt"
	val "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"gitlab.com/distributed_lab/urlval/v4"
	"net/http"
)

type UserInputs struct {
	UserId        string `url:"user_id"`
	AgeLowerBound int    `url:"age_lower_bound"`
	Uniqueness    bool   `url:"uniqueness"`
}

func NewGetUserInputs(r *http.Request) (userInputs UserInputs, err error) {
	if err = urlval.Decode(r.URL.Query(), &userInputs); err != nil {
		err = newDecodeError("query", err)
		return
	}
	err = val.Errors{
		"user_id":         val.Validate(userInputs.UserId, val.Required, is.Email),
		"age_lower_bound": val.Validate(userInputs.AgeLowerBound, val.Required),
		"uniqueness":      val.Validate(userInputs.Uniqueness, val.Required),
	}.
		Filter()
	return
}

func newDecodeError(what string, err error) error {
	return val.Errors{
		what: fmt.Errorf("decode request %s: %w", what, err),
	}
}
