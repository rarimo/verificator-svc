package responses

import (
	"github.com/rarimo/verificator-svc/internal/data"
	"github.com/rarimo/verificator-svc/resources"
)

func NewProofParamsByIdResponse(user data.VerifyUsers, attr resources.ProofParamsAttributes) resources.ProofParamsResponse {
	return resources.ProofParamsResponse{
		Data: resources.ProofParams{
			Key: resources.Key{
				ID:   user.UserIDHash,
				Type: resources.GET_PROOF_PARAMS,
			},
			Attributes: attr,
		},
	}
}
