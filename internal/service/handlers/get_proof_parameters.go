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
	proofSelectorValue = "23073"
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
		Log(r).WithError(err).Errorf("User already exists with userID [%s]", userIdHash)
		ape.RenderErr(w, problems.InternalError())
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
				BirthDateLowerBound:       "18",
				BirthDateUpperBound:       "24",
				CallbackUrl:               fmt.Sprintf("http://localhost:8000/integrations/verificator-svc/public/callback/%s", user.UserIdHash),
				CitizenshipMask:           "city_mask",
				EventData:                 user.UserIdHash,
				EventId:                   "event_id",
				ExpirationDateLowerBound:  "ExpirationDateLowerBound",
				ExpirationDateUpperBound:  "ExpirationDateUpperBound",
				IdentityCounter:           maxIdentityCount,
				IdentityCounterLowerBound: 0,
				IdentityCounterUpperBound: maxIdentityCount,
				Selector:                  proofSelectorValue,
				TimestampLowerBound:       "time_when_app_started",
				TimestampUpperBound:       "time_when_app",
			},
		},
	}
}

func StringToPoseidonHash(inputString string) (string, error) {
	inputBytes, err := hex.DecodeString(inputString)
	if err != nil {
		return "", fmt.Errorf("failed to decode input string: %s", err)
	}
	hash, err := poseidon.HashBytes(inputBytes)
	if err != nil {
		return "", fmt.Errorf("failde to convert input bytes to hash: %s", err)

	}
	return hex.EncodeToString(hash.Bytes()), nil
}
