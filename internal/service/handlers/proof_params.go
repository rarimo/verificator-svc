package handlers

import (
	"github.com/rarimo/verificator-svc/internal/data"
	"github.com/rarimo/verificator-svc/internal/service/requests"
	"github.com/rarimo/verificator-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
	"strconv"
)

type GetProofParam struct {
	eventID                   string
	proofSelector             string
	citizenshipMask           string
	birthDateLowerBound       string
	birthDateUpperBound       string
	timestampUpperBound       string
	timestampLowerBound       string
	identityCounterUpperBound int32
	expirationDateUpperBound  string
	expirationDateLowerBound  string
}

func GetProofParamsById(w http.ResponseWriter, r *http.Request) {
	userIDHash, err := requests.GetProofParamsByID(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	existingUser, err := VerifyUsersQ(r).WhereHashID(userIDHash).Get()
	if err != nil {
		Log(r).WithError(err).Errorf("failed to query user with userID [%s]", userIDHash)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	var (
		eventID                   = ProofParameters(r).EventID
		TimestampUpperBound       = "0"
		IdentityCounterUpperBound int32
	)

	if existingUser.EventId != "" {
		eventID = existingUser.EventId
	}

	birthDateUpperBound := CalculateBirthDateHex(existingUser.AgeLowerBound)
	if existingUser.AgeLowerBound == 0 {
		birthDateUpperBound = "0x303030303030"
	}

	proofSelector := CalculateProofSelector(existingUser.Uniqueness)
	if proofSelector&(1<<timestampUpperBoundBit) != 0 &&
		proofSelector&(1<<identityCounterUpperBoundBit) != 0 {
		TimestampUpperBound = ProofParameters(r).TimestampUpperBound
		IdentityCounterUpperBound = 1
	}

	proofParams := GetProofParam{
		eventID:                   eventID,
		proofSelector:             strconv.Itoa(proofSelector),
		identityCounterUpperBound: IdentityCounterUpperBound,
		timestampUpperBound:       TimestampUpperBound,
		citizenshipMask:           Utf8ToHex(existingUser.Nationality),
		timestampLowerBound:       "0",
		birthDateLowerBound:       "0x303030303030",
		birthDateUpperBound:       birthDateUpperBound,
		expirationDateUpperBound:  "52983525027888",
		expirationDateLowerBound:  "52983525027888",
	}

	ape.Render(w, NewProofParamsByIdResponse(*existingUser, proofParams))

}

func NewProofParamsByIdResponse(user data.VerifyUsers, params GetProofParam) resources.ProofParamsResponse {
	return resources.ProofParamsResponse{
		Data: resources.ProofParams{
			Key: resources.Key{
				ID:   user.UserIDHash,
				Type: resources.GET_PROOF_PARAMS,
			},
			Attributes: resources.ProofParamsAttributes{
				BirthDateLowerBound:       params.birthDateLowerBound,
				BirthDateUpperBound:       params.birthDateUpperBound,
				CitizenshipMask:           params.citizenshipMask,
				EventData:                 user.UserIDHash,
				EventId:                   params.eventID,
				ExpirationDateLowerBound:  params.expirationDateLowerBound,
				ExpirationDateUpperBound:  params.expirationDateUpperBound,
				IdentityCounter:           0,
				IdentityCounterLowerBound: 0,
				IdentityCounterUpperBound: params.identityCounterUpperBound,
				Selector:                  params.proofSelector,
				TimestampLowerBound:       params.timestampLowerBound,
				TimestampUpperBound:       params.timestampUpperBound,
			},
		},
	}
}
