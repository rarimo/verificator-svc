package handlers

import (
	"database/sql"
	"net/http"
	"time"

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

	user := &data.VerifyUsers{
		UserID:        req.Data.ID,
		UserIDHash:    userIDHash,
		CreatedAt:     time.Now().UTC(),
		Status:        "not_verified",
		Proof:         []byte{},
		AgeLowerBound: -1,

		EventID:                   req.Data.Attributes.EventId,
		ExpirationLowerBound:      req.Data.Attributes.ExpirationDateLowerBound,
		BirthDateLowerBound:       sql.NullString{String: req.Data.Attributes.BirthDateLowerBound, Valid: true},
		BirthDateUpperBound:       sql.NullString{String: req.Data.Attributes.BirthDateUpperBound, Valid: true},
		CitizenshipMask:           sql.NullString{String: req.Data.Attributes.CitizenshipMask, Valid: true},
		EventData:                 sql.NullString{String: req.Data.Attributes.EventData, Valid: true},
		ExpirationDateUpperBound:  sql.NullString{String: req.Data.Attributes.ExpirationDateUpperBound, Valid: true},
		IdentityCounter:           req.Data.Attributes.IdentityCounter,
		IdentityCounterLowerBound: req.Data.Attributes.IdentityCounterLowerBound,
		IdentityCounterUpperBound: req.Data.Attributes.IdentityCounterUpperBound,
		Selector:                  req.Data.Attributes.Selector,
		TimestampLowerBound:       sql.NullString{String: req.Data.Attributes.TimestampLowerBound, Valid: true},
		TimestampUpperBound:       sql.NullString{String: req.Data.Attributes.TimestampUpperBound, Valid: true},
	}

	dbUser, err := ctx.VerifyUsersQ(r).Upsert(user)
	if err != nil {
		ctx.Log(r).WithError(err).WithField("user", user).Errorf("failed to upsert user with userID [%s]", user.UserIDHash)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, responses.NewVerificationLinkResponse(dbUser, ctx.Callback(r).URL))
}
