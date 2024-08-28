package requests

import (
	"encoding/json"
	val "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/verificator-svc/resources"

	"net/http"
	"strings"
)

func GetVerificationCallbackByID(r *http.Request) (req resources.ProofRequest, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, newDecodeError("body", err)
	}

	req.Data.ID = strings.ToLower(req.Data.ID)
	var (
		attr = req.Data.Attributes
	)
	return req, val.Errors{
		"data/id":                           val.Validate(req.Data.ID, val.Required),
		"data/type":                         val.Validate(req.Data.Type, val.Required, val.In(resources.RECEIVE_PROOF)),
		"data/attributes/proof":             val.Validate(attr.Proof, val.When(attr.Proof != nil, val.Required)),
		"data/attributes/proof/proof":       val.Validate(attr.Proof.Proof, val.When(attr.Proof.Proof != nil, val.Required)),
		"data/attributes/proof/pub_signals": val.Validate(attr.Proof.PubSignals, val.When(attr.Proof.PubSignals != nil, val.Required, val.Length(22, 22))),
	}.Filter()

}
