package requests

import (
	"encoding/json"
	val "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	zkptypes "github.com/iden3/go-rapidsnark/types"
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
		attr  = req.Data.Attributes
		proof zkptypes.ZKProof
	)
	return req, val.Errors{
		"data/id":                           val.Validate(req.Data.ID, val.Required, is.Hexadecimal),
		"data/attributes/proof":             val.Validate(attr.Proof, val.Required),
		"data/attributes/proof/proof":       val.Validate(proof.Proof, val.When(proof.Proof != nil, val.Required)),
		"data/attributes/proof/pub_signals": val.Validate(proof.PubSignals, val.When(proof.Proof != nil, val.Required, val.Length(22, 22))),
	}.Filter()

}
