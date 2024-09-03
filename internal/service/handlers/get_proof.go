package handlers

import (
	"encoding/json"
	"github.com/iden3/go-rapidsnark/types"
	"github.com/rarimo/verificator-svc/internal/data"
	"github.com/rarimo/verificator-svc/internal/service/requests"
	"github.com/rarimo/verificator-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

func GetProofByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := requests.GetProofByUserID(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	verifiedUser, err := VerifyUsersQ(r).WhereID(userID).Get()
	if err != nil {
		Log(r).WithError(err).Error("failed to get user by userID")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if verifiedUser == nil {
		Log(r).Debugf("User with userID=%s not found", userID)
		ape.RenderErr(w, problems.NotFound())
		return
	}

	ape.Render(w, NewProofByUserIDResponse(*verifiedUser))

}

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
