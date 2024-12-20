package requests

import (
	"encoding/json"
	"net/http"
	"strings"

	val "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/verificator-svc/internal/service/ctx"
	"github.com/rarimo/verificator-svc/resources"
)

func VerificationLink(r *http.Request) (req resources.UserRequest, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, newDecodeError("body", err)
	}

	if !ctx.Verifiers(r).PreserveUserIDCase {
		req.Data.ID = strings.ToLower(req.Data.ID)
	}

	return req, val.Errors{
		"data/id":   val.Validate(req.Data.ID, val.Required),
		"data/type": val.Validate(req.Data.Type, val.Required, val.In(resources.USER)),
	}.Filter()
}
