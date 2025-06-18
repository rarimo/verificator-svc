package requests

import (
	"encoding/json"
	"net/http"
	"strings"

	val "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/verificator-svc/internal/service/ctx"
	"github.com/rarimo/verificator-svc/resources"
)

func VerificationLinkV2(r *http.Request) (req resources.AdvancedVerificationRequest, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, newDecodeError("body", err)
	}

	if !ctx.Verifiers(r).PreserveUserIDCase {
		req.Data.ID = strings.ToLower(req.Data.ID)
	}

	return req, val.Errors{
		"data/id":   val.Validate(req.Data.ID, val.Required),
		"data/type": val.Validate(req.Data.Type, val.Required, val.In(resources.ADVANCED_VERIFICATION)),
		// required
		"data/attributes/event_id": val.Validate(req.Data.Attributes.EventId, val.Required,
			val.NewStringRule(validateEventID, "must be decimal and less than 31 bytes"),
		),
		"data/attributes/selector": val.Validate(req.Data.Attributes.Selector, val.Required),
		// base v1
		"data/attributes/age_lower_bound":        val.Validate(req.Data.Attributes.AgeLowerBound, val.NilOrNotEmpty),
		"data/attributes/citizenship_mask":       val.Validate(req.Data.Attributes.CitizenshipMask, val.NilOrNotEmpty),
		"data/attributes/expiration_lower_bound": val.Validate(req.Data.Attributes.ExpirationLowerBound, val.NilOrNotEmpty),
		// advanced v2
		"data/attributes/identity_counter":             val.Validate(req.Data.Attributes.IdentityCounter, val.When(req.Data.Attributes.IdentityCounter != nil, val.Min(0))),
		"data/attributes/identity_counter_lower_bound": val.Validate(req.Data.Attributes.IdentityCounterLowerBound, val.When(req.Data.Attributes.IdentityCounterLowerBound != nil, val.Min(0))),
		"data/attributes/identity_counter_upper_bound": val.Validate(req.Data.Attributes.IdentityCounterUpperBound, val.When(req.Data.Attributes.IdentityCounterUpperBound != nil, val.Min(0))),
		"data/attributes/birth_date_lower_bound":       val.Validate(req.Data.Attributes.BirthDateLowerBound, val.NilOrNotEmpty),
		"data/attributes/birth_date_upper_bound":       val.Validate(req.Data.Attributes.BirthDateUpperBound, val.NilOrNotEmpty),
		"data/attributes/event_data":                   val.Validate(req.Data.Attributes.EventData, val.NilOrNotEmpty),
		"data/attributes/expiration_date_lower_bound":  val.Validate(req.Data.Attributes.ExpirationDateLowerBound, val.NilOrNotEmpty),
		"data/attributes/expiration_date_upper_bound":  val.Validate(req.Data.Attributes.ExpirationDateUpperBound, val.NilOrNotEmpty),
		"data/attributes/timestamp_lower_bound":        val.Validate(req.Data.Attributes.TimestampLowerBound, val.When(req.Data.Attributes.TimestampLowerBound != nil, val.Min(0))),
		"data/attributes/timestamp_upper_bound":        val.Validate(req.Data.Attributes.TimestampUpperBound, val.When(req.Data.Attributes.TimestampUpperBound != nil, val.Min(0))),
	}.Filter()
}
