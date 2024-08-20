package handlers

import (
	"encoding/hex"
	"fmt"
	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/rarimo/verificator-svc/internal/data"
	"github.com/rarimo/verificator-svc/internal/service/requests"
	"github.com/rarimo/verificator-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
	"time"
)

const (
	maxIdentityCount   = 1
	proofSelectorValue = "236065"
)

func GetProofParameters(w http.ResponseWriter, r *http.Request) {
	userInputs, err := requests.NewGetUserInputs(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	userIdHash, err := StringToPoseidonHash(userInputs.UserId)
	if err != nil {
		ape.RenderErr(w, problems.InternalError())
		return
	}
	user := &data.VerifyUsers{
		UserID:     userInputs.UserId,
		UserIdHash: userIdHash,
		CreatedAt:  time.Now().UTC(),
		Status:     "false",
	}

	existingUser, err := VerifyUsersQ(r).WhereHashID(user.UserIdHash).Get()
	if err != nil {
		Log(r).WithError(err).Errorf("failed to query user with userID [%s]", userIdHash)
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if existingUser != nil {
		ape.Render(w, NewProofParametersResponse(*existingUser))
		return
	}

	err = VerifyUsersQ(r).Insert(user)
	if err != nil {
		Log(r).WithError(err).Errorf("failed to insert user with userID [%s]", user.UserIdHash)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, NewProofParametersResponse(*user))
}

func NewProofParametersResponse(user data.VerifyUsers) resources.ParametersResponse {
	return resources.ParametersResponse{
		Data: resources.Parameters{
			Key: resources.Key{
				ID:   user.UserID,
				Type: resources.PROOF_PARAMETERS,
			},
			Attributes: resources.ParametersAttributes{
				BirthDateLowerBound:       "313031303130",
				BirthDateUpperBound:       "323031303130",
				CallbackUrl:               fmt.Sprintf("http://localhost:8000/integrations/verificator-svc/public/callback/%s", user.UserIdHash),
				CitizenshipMask:           "0",
				EventData:                 user.UserIdHash,
				EventId:                   "111186066134341633902189494613533900917417361106374681011849132651019822199",
				ExpirationDateLowerBound:  "313031303130",
				ExpirationDateUpperBound:  "333031303130",
				IdentityCounter:           0,
				IdentityCounterLowerBound: 0,
				IdentityCounterUpperBound: maxIdentityCount,
				Selector:                  proofSelectorValue,
				TimestampLowerBound:       "10000000000",
				TimestampUpperBound:       "19000000000",
			},
		},
	}
}

func StringToPoseidonHash(inputString string) (string, error) {
	inputBytes := []byte(inputString)

	hash, err := poseidon.HashBytes(inputBytes)
	if err != nil {
		return "", fmt.Errorf("failde to convert input bytes to hash: %s", err)

	}
	return hex.EncodeToString(hash.Bytes()), nil
}
