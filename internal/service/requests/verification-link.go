package requests

import (
	"encoding/json"
	val "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/verificator-svc/resources"
	"net/http"
	"strings"
)

func VerificationLink(r *http.Request) (req resources.UserRequest, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, newDecodeError("body", err)
	}

	req.Data.ID = strings.ToLower(req.Data.ID)
	var (
		attr = req.Data.Attributes
	)

	return req, val.Errors{
		"data/id":                         val.Validate(req.Data.ID, val.Required),
		"data/type":                       val.Validate(req.Data.Type, val.Required, val.In(resources.USER)),
		"data/attributes/age_lower_bound": val.Validate(attr.AgeLowerBound, val.Required),
		"data/attributes/uniqueness":      val.Validate(attr.Uniqueness, val.Required),
		"data/attributes/nationality":     val.Validate(attr.Nationality, val.Required),
		"data/attributes/event_id":        val.Validate(attr.Nationality, val.Required),
	}.Filter()
}
