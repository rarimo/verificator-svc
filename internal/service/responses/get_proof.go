package responses

import (
	"encoding/json"
	"github.com/iden3/go-rapidsnark/types"
	"github.com/rarimo/verificator-svc/internal/data"
	"github.com/rarimo/verificator-svc/resources"
)

func NewProofByUserIDResponse(user data.VerifyUsers) resources.GetProofRequest {
	var proof types.ZKProof
	_ = json.Unmarshal(user.Proof, &proof)
	return resources.GetProofRequest{
		Data: resources.GetProof{
			Key: resources.Key{
				ID:   user.UserIDHash,
				Type: resources.GET_PROOF,
			},
			Attributes: resources.GetProofAttributes{
				Proof: proof,
			},
		},
	}
}
