package handlers

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userIdHash := md5.Sum([]byte(userInputs.UserId))
	user := &data.VerifyUsers{
		UserID:     userInputs.UserId,
		UserIdHash: hex.EncodeToString(userIdHash[:]),
		CreatedAt:  time.Now().UTC(),
		Status:     "false",
	}

	existingUser, err := VerifyUsersQ(r).WhereHashID(user.UserIdHash).Get()
	if err != nil {
		Log(r).WithError(err).Errorf("failed to query user with userID [%s]", hex.EncodeToString(userIdHash[:]))
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if existingUser != nil {
		Log(r).WithError(err).Errorf("User already exists with userID [%s]", hex.EncodeToString(userIdHash[:]))
		ape.RenderErr(w, problems.InternalError())
		return
	}

	err = VerifyUsersQ(r).Insert(user)
	if err != nil {
		Log(r).WithError(err).Errorf("failed to insert user with userID [%s]", user.UserIdHash)
		ape.RenderErr(w, problems.InternalError())
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
