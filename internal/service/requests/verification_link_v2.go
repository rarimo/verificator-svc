package requests

import (
	"encoding/json"
	"net/http"
	"strings"

	val "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/verificator-svc/internal/service/ctx"
	"github.com/rarimo/verificator-svc/resources"
)

func VerificationLinkV2(r *http.Request) (req resources.UserV2Request, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, newDecodeError("body", err)
	}

	if !ctx.Verifiers(r).PreserveUserIDCase {
		req.Data.ID = strings.ToLower(req.Data.ID)
	}

	return req, val.Errors{
		"data/id":                                      val.Validate(req.Data.ID, val.Required),
		"data/type":                                    val.Validate(req.Data.Type, val.Required, val.In(resources.USER_V2)),
		"data/attributes/birth_date_lower_bound":       val.Validate(req.Data.Attributes.BirthDateLowerBound, val.Required),
		"data/attributes/birth_date_upper_bound":       val.Validate(req.Data.Attributes.BirthDateUpperBound, val.Required),
		"data/attributes/citizenship_mask":             val.Validate(req.Data.Attributes.CitizenshipMask, val.Required),
		"data/attributes/event_data":                   val.Validate(req.Data.Attributes.EventData, val.Required),
		"data/attributes/event_id":                     val.Validate(req.Data.Attributes.EventId, val.Required),
		"data/attributes/expiration_date_lower_bound":  val.Validate(req.Data.Attributes.ExpirationDateLowerBound, val.Required),
		"data/attributes/expiration_date_upper_bound":  val.Validate(req.Data.Attributes.ExpirationDateUpperBound, val.Required),
		"data/attributes/identity_counter":             val.Validate(req.Data.Attributes.IdentityCounter, val.Required),
		"data/attributes/identity_counter_lower_bound": val.Validate(req.Data.Attributes.IdentityCounterLowerBound, val.Required),
		"data/attributes/identity_counter_upper_bound": val.Validate(req.Data.Attributes.IdentityCounterUpperBound, val.Required),
		"data/attributes/selector":                     val.Validate(req.Data.Attributes.Selector, val.Required),
		"data/attributes/timestamp_lower_bound":        val.Validate(req.Data.Attributes.TimestampLowerBound, val.Required),
		"data/attributes/timestamp_upper_bound":        val.Validate(req.Data.Attributes.TimestampUpperBound, val.Required),
	}.Filter()
}
