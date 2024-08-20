package requests

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
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
	return
}

func newDecodeError(what string, err error) error {
	return validation.Errors{
		what: fmt.Errorf("decode request %s: %w", what, err),
	}
}
