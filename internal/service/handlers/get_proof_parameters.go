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

const (
	nullifierBit                 = 0
	citizenshipBit               = 5
	timestampUpperBoundBit       = 9
	identityCounterUpperBoundBit = 11
	expirationDateLowerboundBit  = 12
	expirationDateUpperbound     = 13
	birthDateLowerboundBit       = 14
	birthDateUpperboundBit       = 15
	birthDateFormat              = "060102"
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
		return
	}

	var (
		eventID                   = ProofParameters(r).EventID
		TimestampUpperBound       = "0"
		IdentityCounterUpperBound int32
	)

	if userInputs.EventID != "" {
		eventID = userInputs.EventID
	}

	proofSelector := CalculateProofSelector(userInputs.Uniqueness)
	if proofSelector&(1<<timestampUpperBoundBit) != 0 &&
		proofSelector&(1<<identityCounterUpperBoundBit) != 0 {
		TimestampUpperBound = ProofParameters(r).TimestampUpperBound
		IdentityCounterUpperBound = 1
	}

	userIdHash, err := StringToPoseidonHash(userInputs.UserId)
	if err != nil {
		Log(r).WithError(err).Errorf("failed to convert user with userID [%s] to poseidon hash", userInputs.UserId)
		ape.RenderErr(w, problems.InternalError())
		return
	}
	user := &data.VerifyUsers{
		UserID:        userInputs.UserId,
		UserIDHash:    userIdHash,
		CreatedAt:     time.Now().UTC(),
		Status:        "not_verified",
		Nationality:   userInputs.Nationality,
		AgeLowerBound: userInputs.AgeLowerBound,
		Uniqueness:    userInputs.Uniqueness,
		Proof:         []byte{},
	}

	proofParams := ProofParams{
		host:                      Callback(r).URL,
		eventID:                   eventID,
		proofSelector:             strconv.Itoa(proofSelector),
		identityCounterUpperBound: IdentityCounterUpperBound,
		timestampUpperBound:       TimestampUpperBound,
		citizenshipMask:           Utf8ToHex(userInputs.Nationality),
		timestampLowerBound:       "0",
		birthDateLowerBound:       "0x303030303030",
		birthDateUpperBound:       CalculateBirthDateHex(userInputs.AgeLowerBound),
		expirationDateUpperBound:  "52983525027888",
		expirationDateLowerBound:  "52983525027888",
	}

	existingUser, err := VerifyUsersQ(r).WhereHashID(user.UserIDHash).Get()
	if err != nil {
		Log(r).WithError(err).Errorf("failed to query user with userID [%s]", userIdHash)
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if existingUser != nil {
		ape.Render(w, NewProofParametersResponse(*existingUser, proofParams))
		return
	}

	if err = VerifyUsersQ(r).Insert(user); err != nil {
		Log(r).WithError(err).Errorf("failed to insert user with userID [%s]", user.UserIDHash)
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
				CallbackUrl:               fmt.Sprintf("%s/integrations/verificator-svc/public/callback/%s", params.host, user.UserIDHash),
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

func StringToPoseidonHash(inputString string) (string, error) {
	inputBytes := []byte(inputString)

	hash, err := poseidon.HashBytes(inputBytes)
	if err != nil {
		return "", fmt.Errorf("failde to convert input bytes to hash: %w", err)

	}
	return fmt.Sprintf("0x%s", hex.EncodeToString(hash.Bytes())), nil
}

func Utf8ToHex(input string) string {
	bytes := []byte(input)
	hexString := hexutils.BytesToHex(bytes)
	return fmt.Sprintf("0x%s", hexString)
}

func CalculateBirthDateHex(ageLowerBound int) string {
	allowedBirthDate := time.Now().UTC().AddDate(-ageLowerBound, 0, 0)
	formattedDate := []byte(allowedBirthDate.Format(birthDateFormat))
	hexBirthDateLoweBound := hexutils.BytesToHex(formattedDate)

	return fmt.Sprintf("0x%s", hexBirthDateLoweBound)
}

func CalculateProofSelector(uniqueness bool) int {
	var bitLine uint32
	bitLine |= 1 << nullifierBit
	bitLine |= 1 << citizenshipBit
	bitLine |= 1 << birthDateUpperboundBit
	if uniqueness {
		bitLine |= 1 << timestampUpperBoundBit
		bitLine |= 1 << identityCounterUpperBoundBit
	}

	return int(bitLine)
}
