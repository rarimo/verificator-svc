package requests

import (
	"encoding/json"
	"math/big"
	"net/http"
	"strings"

	val "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/verificator-svc/internal/config"
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
		"data/attributes/event_id": val.Validate(req.Data.Attributes.EventId, val.NilOrNotEmpty,
			val.When(
				!val.IsEmpty(req.Data.Attributes.EventId),
				val.NewStringRule(validateEventID, "must be decimal and less than 254 bits"),
			),
		),
	}.Filter()
}

func validateEventID(value string) bool {
	eventID, ok := new(big.Int).SetString(value, 10)
	if !ok {
		return false
	}

	if eventID.Cmp(big.NewInt(0)) <= 0 {
		return false
	}

	if eventID.Cmp(config.MaxEventId) == 1 {
		return false
	}

	return true
}
