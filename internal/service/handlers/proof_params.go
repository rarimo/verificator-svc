package handlers

import (
	"fmt"
	"github.com/rarimo/verificator-svc/internal/service/handlers/helpers"
	"github.com/rarimo/verificator-svc/internal/service/requests"
	"github.com/rarimo/verificator-svc/internal/service/responses"
	"github.com/rarimo/verificator-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
	"strconv"
)

func GetProofParamsById(w http.ResponseWriter, r *http.Request) {
	userIDHash, err := requests.GetProofParamsByID(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	existingUser, err := VerifyUsersQ(r).WhereHashID(userIDHash).Get()
	if err != nil {
		Log(r).WithError(err).Errorf("failed to query user with userID [%s]", userIDHash)
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	if existingUser == nil {
		Log(r).WithError(err).Errorf("user with userID [%s] not found", userIDHash)
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	var (
		IdentityCounterUpperBound int32
		TimestampUpperBound       = "0"
		eventID                   = Verifiers(r).EventID
		birthDateUpperBound       = helpers.CalculateBirthDateHex(existingUser.AgeLowerBound)
		proofSelector             = helpers.CalculateProofSelector(existingUser.Uniqueness, existingUser.AgeLowerBound, existingUser.Nationality, existingUser.SexEnable, existingUser.NationalityEnable)
		callbackURL               = fmt.Sprintf("%s/integrations/verificator-svc/public/callback/%s", Callback(r).URL, userIDHash)
	)

	if existingUser.EventId != "" {
		eventID = existingUser.EventId
	}

	if existingUser.AgeLowerBound == -1 {
		birthDateUpperBound = "0x303030303030"
	}

	if proofSelector&(1<<helpers.TimestampUpperBoundBit) != 0 &&
		proofSelector&(1<<helpers.IdentityCounterUpperBoundBit) != 0 {
		TimestampUpperBound = strconv.FormatInt(Verifiers(r).ServiceStartTimestamp, 10)
		IdentityCounterUpperBound = 1
	}

	proofParameters := resources.ProofParamsAttributes{
		BirthDateLowerBound:       "0x303030303030",
		BirthDateUpperBound:       birthDateUpperBound,
		CitizenshipMask:           helpers.Utf8ToHex(existingUser.Nationality),
		EventData:                 existingUser.UserIDHash,
		EventId:                   eventID,
		ExpirationDateLowerBound:  "52983525027888",
		ExpirationDateUpperBound:  "52983525027888",
		IdentityCounter:           0,
		IdentityCounterLowerBound: 0,
		IdentityCounterUpperBound: IdentityCounterUpperBound,
		Selector:                  strconv.Itoa(proofSelector),
		TimestampLowerBound:       "0",
		TimestampUpperBound:       TimestampUpperBound,
		CallbackUrl:               &callbackURL,
	}

	ape.Render(w, responses.NewProofParamsByIdResponse(*existingUser, proofParameters))

}
