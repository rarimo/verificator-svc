package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/verificator-svc/internal/data"
	"github.com/rarimo/verificator-svc/internal/service/ctx"
	"github.com/rarimo/verificator-svc/internal/service/handlers/helpers"
	"github.com/rarimo/verificator-svc/internal/service/requests"
	"github.com/rarimo/verificator-svc/internal/service/responses"
	"github.com/rarimo/web3-auth-svc/pkg/auth"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func VerificationLinkV2(w http.ResponseWriter, r *http.Request) {
	req, err := requests.VerificationLinkV2(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if !helpers.Authenticates(ctx.AuthClient(r), ctx.UserClaims(r), auth.UserGrant(req.Data.ID)) {
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	userIDHash, err := helpers.BuildUserIDHash(req.Data.ID, ctx.Verifiers(r).ERC1155)
	if err != nil {
		ctx.Log(r).WithError(err).WithField("user_id", req.Data.ID).Error("error building user ID hash")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	selectorInt64, err := strconv.ParseInt(req.Data.Attributes.Selector, 10, 32)
	if err != nil {
		ctx.Log(r).WithError(err).WithField("selector", req.Data.Attributes.Selector).Error("failed to parse selector")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"selector": err,
		})...)
		return
	}

	user := &data.VerifyUsers{
		UserID:               req.Data.ID,
		UserIDHash:           userIDHash,
		CreatedAt:            time.Now().UTC(),
		Status:               "not_verified",
		Proof:                []byte{},
		AgeLowerBound:        -1,
		ExpirationLowerBound: helpers.DefaultDateHex,

		// required v2
		EventID:                   req.Data.Attributes.EventId,
		Selector:                  int32(selectorInt64),
		IdentityCounter:           0,
		IdentityCounterLowerBound: 0,
		IdentityCounterUpperBound: 0,
	}

	// base v1
	if req.Data.Attributes.CitizenshipMask != nil {
		user.Nationality = *req.Data.Attributes.CitizenshipMask
	}

	if req.Data.Attributes.AgeLowerBound != nil {
		user.AgeLowerBound = int(*req.Data.Attributes.AgeLowerBound)
	}

	if req.Data.Attributes.ExpirationLowerBound != nil {
		user.ExpirationLowerBound = helpers.GetExpirationLowerBound(*req.Data.Attributes.ExpirationLowerBound)
	}

	// V2 identity counter
	if req.Data.Attributes.IdentityCounter != nil {
		user.IdentityCounter = *req.Data.Attributes.IdentityCounter
	}

	if req.Data.Attributes.IdentityCounterLowerBound != nil {
		user.IdentityCounterLowerBound = *req.Data.Attributes.IdentityCounterLowerBound
	}

	if req.Data.Attributes.IdentityCounterUpperBound != nil {
		user.IdentityCounterUpperBound = *req.Data.Attributes.IdentityCounterUpperBound
	}

	// V2 nullable
	if req.Data.Attributes.BirthDateLowerBound != nil {
		user.BirthDateLowerBound = sql.NullString{String: *req.Data.Attributes.BirthDateLowerBound, Valid: true}
	}

	if req.Data.Attributes.BirthDateUpperBound != nil {
		user.BirthDateUpperBound = sql.NullString{String: *req.Data.Attributes.BirthDateUpperBound, Valid: true}
	}

	if req.Data.Attributes.EventData != nil {
		eventDataDecimal, err := helpers.HexToDecimal(*req.Data.Attributes.EventData)
		if err != nil {
			ctx.Log(r).WithError(err).Error("failed to convert event_data from hex to decimal")
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"event_data": err,
			})...)
			return
		}
		user.EventData = sql.NullString{String: eventDataDecimal, Valid: true}
	}

	if req.Data.Attributes.ExpirationDateUpperBound != nil {
		user.ExpirationDateUpperBound = sql.NullString{String: *req.Data.Attributes.ExpirationDateUpperBound, Valid: true}
	}

	if req.Data.Attributes.TimestampLowerBound != nil {
		user.TimestampLowerBound = sql.NullTime{Time: time.Unix(*req.Data.Attributes.TimestampLowerBound, 0), Valid: true}
	}

	if req.Data.Attributes.TimestampUpperBound != nil {
		user.TimestampUpperBound = sql.NullTime{Time: time.Unix(*req.Data.Attributes.TimestampUpperBound, 0), Valid: true}
	}

	if req.Data.Attributes.Sex != nil {
		user.Sex = *req.Data.Attributes.Sex
	}

	dbUser, err := ctx.VerifyUsersQ(r).Upsert(user)
	if err != nil {
		ctx.Log(r).WithError(err).WithField("user", user).Errorf("failed to upsert user with userID [%s]", user.UserIDHash)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, responses.NewVerificationLinkResponse(dbUser, ctx.Callback(r).URL))
}
