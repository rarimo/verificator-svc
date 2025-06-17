package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/rarimo/verificator-svc/internal/service/ctx"
	"github.com/rarimo/verificator-svc/internal/service/handlers/helpers"
	"github.com/rarimo/verificator-svc/internal/service/requests"
	"github.com/rarimo/verificator-svc/internal/service/responses"
	"github.com/rarimo/verificator-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetProofParamsById(w http.ResponseWriter, r *http.Request) {
	userIDHash, err := requests.GetPathUserIDHash(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	existingUser, err := ctx.VerifyUsersQ(r).WhereHashID(userIDHash).Get()
	if err != nil {
		ctx.Log(r).WithError(err).Errorf("failed to query user with userID [%s]", userIDHash)
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	if existingUser == nil {
		ctx.Log(r).WithError(err).Errorf("user with userID [%s] not found", userIDHash)
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	callbackURL := fmt.Sprintf("%s/integrations/verificator-svc/public/callback/%s", ctx.Callback(r).URL, userIDHash)

	// check if selector is not empty to use v2 params
	if existingUser.Selector != "" {
		proofParameters := helpers.BuildV2ProofParams(existingUser, callbackURL)
		ape.Render(w, responses.NewProofParamsByIdResponse(*existingUser, proofParameters))
		return
	}

	var (
		IdentityCounterUpperBound int32
		TimestampUpperBound       = "0"
		eventID                   = ctx.Verifiers(r).EventID
		birthDateUpperBound       = helpers.CalculateBirthDateHex(existingUser.AgeLowerBound)
		proofSelector             = helpers.CalculateProofSelector(helpers.SelectorParams{
			Uniqueness:           existingUser.Uniqueness,
			AgeLowerBound:        existingUser.AgeLowerBound,
			Nationality:          existingUser.Nationality,
			SexEnable:            existingUser.SexEnable,
			NationalityEnable:    existingUser.NationalityEnable,
			ExpirationLowerBound: !helpers.IsDefaultZKDate(existingUser.ExpirationLowerBound), // If there is non-default value, selector should be enabled
		})
	)

	if existingUser.EventID != "" {
		eventID = existingUser.EventID
	}

	if existingUser.AgeLowerBound == -1 {
		birthDateUpperBound = helpers.DefaultDateHex
	}

	if proofSelector&(1<<helpers.TimestampUpperBoundBit) != 0 &&
		proofSelector&(1<<helpers.IdentityCounterUpperBoundBit) != 0 {
		TimestampUpperBound = strconv.FormatInt(ctx.Verifiers(r).ServiceStartTimestamp, 10)
		IdentityCounterUpperBound = 1
	}

	proofParameters := resources.ProofParamsAttributes{
		BirthDateLowerBound:       helpers.DefaultDateHex,
		BirthDateUpperBound:       birthDateUpperBound,
		CitizenshipMask:           helpers.Utf8ToHex(existingUser.Nationality),
		EventData:                 existingUser.UserIDHash,
		EventId:                   eventID,
		ExpirationDateLowerBound:  existingUser.ExpirationLowerBound,
		ExpirationDateUpperBound:  helpers.DefaultDateHex,
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
