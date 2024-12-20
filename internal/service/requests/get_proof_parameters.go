package requests

import (
	"fmt"
	"net/http"
	"strings"

	val "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/rarimo/verificator-svc/internal/service/ctx"
	"gitlab.com/distributed_lab/urlval/v4"
)

type UserInputs struct {
	UserID               string `url:"user_id"`
	AgeLowerBound        int    `url:"age_lower_bound"`
	Uniqueness           bool   `url:"uniqueness"`
	Nationality          string `url:"nationality"`
	EventID              string `url:"event_id"`
	ExpirationLowerBound bool   `url:"expiration_lower_bound"`
}

func NewGetUserInputs(r *http.Request) (userInputs UserInputs, err error) {
	if err = urlval.Decode(r.URL.Query(), &userInputs); err != nil {
		err = newDecodeError("query", err)
		return
	}

	if !ctx.Verifiers(r).PreserveUserIDCase {
		userInputs.UserID = strings.ToLower(userInputs.UserID)
	}

	err = val.Errors{
		"user_id":         val.Validate(userInputs.UserID, val.Required),
		"age_lower_bound": val.Validate(userInputs.AgeLowerBound, val.Required),
		"uniqueness":      val.Validate(val.Required),
		"nationality":     val.Validate(userInputs.Nationality, val.Required),
		"event_id":        val.Validate(userInputs.EventID, is.Hexadecimal),
	}.Filter()
	return
}

func newDecodeError(what string, err error) error {
	return val.Errors{
		what: fmt.Errorf("decode request %s: %w", what, err),
	}
}
