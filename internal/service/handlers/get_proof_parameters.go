package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rarimo/verificator-svc/internal/data"
	"github.com/rarimo/verificator-svc/internal/service/handlers/helpers"
	"github.com/rarimo/verificator-svc/internal/service/requests"
	"github.com/rarimo/verificator-svc/internal/service/responses"
	"github.com/rarimo/verificator-svc/resources"
	"github.com/rarimo/web3-auth-svc/pkg/auth"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetProofParameters(w http.ResponseWriter, r *http.Request) {
	userInputs, err := requests.NewGetUserInputs(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if !helpers.Authenticates(AuthClient(r), UserClaims(r), auth.UserGrant(userInputs.UserId)) {
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	var (
		IdentityCounterUpperBound int32
		TimestampUpperBound       = "0"
		eventID                   = Verifiers(r).EventID
		proofSelector             = helpers.CalculateProofSelector(helpers.SelectorParams{
			Uniqueness:           userInputs.Uniqueness,
			AgeLowerBound:        userInputs.AgeLowerBound,
			Nationality:          userInputs.Nationality,
			SexEnable:            true,
			NationalityEnable:    true,
			ExpirationLowerBound: userInputs.ExpirationLowerBound,
		})
		expirationLowerBound = helpers.GetExpirationLowerBound(userInputs.ExpirationLowerBound)
	)

	if userInputs.EventID != "" {
		eventID = userInputs.EventID
	}

	if proofSelector&(1<<helpers.TimestampUpperBoundBit) != 0 &&
		proofSelector&(1<<helpers.IdentityCounterUpperBoundBit) != 0 {
		TimestampUpperBound = strconv.FormatInt(Verifiers(r).ServiceStartTimestamp, 10)
		IdentityCounterUpperBound = 1
	}

	userIdHash, err := helpers.StringToPoseidonHash(userInputs.UserId)
	if err != nil {
		Log(r).WithError(err).Errorf("failed to convert user with userID [%s] to poseidon hash", userInputs.UserId)
		ape.RenderErr(w, problems.InternalError())
		return
	}
	user := &data.VerifyUsers{
		UserID:               userInputs.UserId,
		UserIDHash:           userIdHash,
		CreatedAt:            time.Now().UTC(),
		Status:               "not_verified",
		Nationality:          userInputs.Nationality,
		AgeLowerBound:        userInputs.AgeLowerBound,
		Uniqueness:           userInputs.Uniqueness,
		Proof:                []byte{},
		ExpirationLowerBound: expirationLowerBound,
	}

	proofParameters := resources.ParametersAttributes{
		BirthDateLowerBound:       helpers.DefaultDateHex,
		BirthDateUpperBound:       helpers.CalculateBirthDateHex(userInputs.AgeLowerBound),
		CallbackUrl:               fmt.Sprintf("%s/integrations/verificator-svc/public/callback/%s", Callback(r).URL, user.UserIDHash),
		CitizenshipMask:           helpers.Utf8ToHex(userInputs.Nationality),
		EventData:                 helpers.GetEventData(common.HexToAddress(user.UserID).Bytes()),
		EventId:                   eventID,
		ExpirationDateLowerBound:  expirationLowerBound,
		ExpirationDateUpperBound:  helpers.DefaultDateHex,
		IdentityCounter:           0,
		IdentityCounterLowerBound: 0,
		IdentityCounterUpperBound: IdentityCounterUpperBound,
		Selector:                  strconv.Itoa(proofSelector),
		TimestampLowerBound:       "0",
		TimestampUpperBound:       TimestampUpperBound,
	}

	existingUser, err := VerifyUsersQ(r).WhereHashID(user.UserIDHash).Get()
	if err != nil {
		Log(r).WithError(err).Errorf("failed to query user with userID [%s]", user.UserIDHash)
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if existingUser != nil {
		ape.Render(w, responses.NewProofParametersResponse(*existingUser, proofParameters))
		return
	}

	if err = VerifyUsersQ(r).Insert(user); err != nil {
		Log(r).WithError(err).Errorf("failed to insert user with userID [%s]", user.UserIDHash)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, responses.NewProofParametersResponse(*user, proofParameters))
}
