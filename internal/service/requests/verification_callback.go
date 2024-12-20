package requests

import (
	"encoding/json"
	"net/http"
	"strings"

	val "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/verificator-svc/internal/service/ctx"
	"github.com/rarimo/verificator-svc/resources"
)

func GetVerificationCallbackByID(r *http.Request) (req resources.ProofRequest, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, newDecodeError("body", err)
	}

	if ctx.Verifiers(r).LowerCaseUserID {
		req.Data.ID = strings.ToLower(req.Data.ID)
	}
	var attr = req.Data.Attributes

	err = val.Errors{
		"data/id":               val.Validate(req.Data.ID, val.Required),
		"data/type":             val.Validate(req.Data.Type, val.Required, val.In(resources.RECEIVE_PROOF)),
		"data/attributes/proof": val.Validate(attr.Proof, val.Required),
	}.Filter()

	if err != nil {
		return req, err
	}

	return req, val.Errors{
		"data/attributes/proof/proof":       val.Validate(attr.Proof.Proof, val.Required),
		"data/attributes/proof/pub_signals": val.Validate(attr.Proof.PubSignals, val.Required, val.Length(23, 23)),
	}.Filter()

}
