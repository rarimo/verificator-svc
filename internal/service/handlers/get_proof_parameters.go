package handlers

import (
	"encoding/hex"
	"fmt"
	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/rarimo/verificator-svc/internal/data"
	"github.com/rarimo/verificator-svc/internal/service/requests"
	"github.com/rarimo/verificator-svc/resources"
	"github.com/status-im/keycard-go/hexutils"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
	"strconv"
	"time"
)

type ProofParams struct {
	host                      string
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

func GetProofParameters(w http.ResponseWriter, r *http.Request) {
	userInputs, err := requests.NewGetUserInputs(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		Log(r).Debug(userInputs.Uniqueness)
		return
	}

	TimestampUpperBound := "0"
	var IdentityCounterUpperBound int32 = 0
	proofSelector := CalculateProofSelector(userInputs.Uniqueness)
	if proofSelector&(1<<9) != 0 && proofSelector&(1<<11) != 0 {
		TimestampUpperBound = ProofParameters(r).TimestampUpperBound
		IdentityCounterUpperBound = 1
	}

	userIdHash, err := StringToPoseidonHash(userInputs.UserId)
	if err != nil {
		ape.RenderErr(w, problems.InternalError())
		return
	}
	user := &data.VerifyUsers{
		UserID:        userInputs.UserId,
		UserIdHash:    userIdHash,
		CreatedAt:     time.Now().UTC(),
		Status:        "not_verified",
		Nationality:   userInputs.Nationality,
		AgeLowerBound: userInputs.AgeLowerBound,
		Uniqueness:    userInputs.Uniqueness,
	}

	proofParams := ProofParams{
		host:                      Callback(r).Url,
		eventID:                   ProofParameters(r).EventID,
		proofSelector:             strconv.Itoa(proofSelector),
		identityCounterUpperBound: IdentityCounterUpperBound,
		timestampUpperBound:       TimestampUpperBound,
		citizenshipMask:           utf8ToHex(userInputs.Nationality),
		timestampLowerBound:       ProofParameters(r).TimestampLowerBound,
		expirationDateUpperBound:  ProofParameters(r).ExpirationDateUpperBound,
		expirationDateLowerBound:  ProofParameters(r).ExpirationDateLowerBound,
		birthDateLowerBound:       calculateBirthDateHex(userInputs.AgeLowerBound),
		birthDateUpperBound:       ProofParameters(r).BirthDateUpperBound,
	}

	existingUser, err := VerifyUsersQ(r).WhereHashID(user.UserIdHash).Get()
	if err != nil {
		Log(r).WithError(err).Errorf("failed to query user with userID [%s]", userIdHash)
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if existingUser != nil {
		ape.Render(w, NewProofParametersResponse(*existingUser, proofParams))
		return
	}

	err = VerifyUsersQ(r).Insert(user)
	if err != nil {
		Log(r).WithError(err).Errorf("failed to insert user with userID [%s]", user.UserIdHash)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, NewProofParametersResponse(*user, proofParams))
}

func NewProofParametersResponse(user data.VerifyUsers, params ProofParams) resources.ParametersResponse {
	return resources.ParametersResponse{
		Data: resources.Parameters{
			Key: resources.Key{
				ID:   user.UserID,
				Type: resources.PROOF_PARAMETERS,
			},
			Attributes: resources.ParametersAttributes{
				BirthDateLowerBound:       params.birthDateLowerBound,
				BirthDateUpperBound:       params.birthDateUpperBound,
				CallbackUrl:               fmt.Sprintf("%s/integrations/verificator-svc/public/callback/%s", params.host, user.UserIdHash),
				CitizenshipMask:           params.citizenshipMask,
				EventData:                 user.UserIdHash,
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

func StringToPoseidonHash(inputString string) (string, error) {
	inputBytes := []byte(inputString)

	hash, err := poseidon.HashBytes(inputBytes)
	if err != nil {
		return "", fmt.Errorf("failde to convert input bytes to hash: %s", err)

	}
	return hex.EncodeToString(hash.Bytes()), nil
}

func utf8ToHex(input string) string {
	bytes := []byte(input)
	hexString := hexutils.BytesToHex(bytes)
	return hexString
}

func calculateBirthDateHex(ageLowerBound int) string {
	currentDate := time.Now().UTC()
	birthDateLoweBound := []byte(fmt.Sprintf("%02d", (currentDate.Year()-ageLowerBound)%100) + "0101")
	hexBirthDateLoweBound := hexutils.BytesToHex(birthDateLoweBound)

	return hexBirthDateLoweBound
}

func CalculateProofSelector(uniqueness bool) int {
	var bitLine uint32 = 0
	if uniqueness {
		bitLine |= 1 << 9
		bitLine |= 1 << 11
	}

	return int(bitLine)
}
