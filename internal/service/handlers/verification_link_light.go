package handlers

import (
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rarimo/verificator-svc/internal/data"
	"github.com/rarimo/verificator-svc/internal/service/handlers/helpers"
	"github.com/rarimo/verificator-svc/internal/service/requests"
	"github.com/rarimo/verificator-svc/internal/service/responses"
	"github.com/rarimo/web3-auth-svc/pkg/auth"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func VerificationLinkLight(w http.ResponseWriter, r *http.Request) {
	req, err := requests.VerificationLink(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if !helpers.Authenticates(AuthClient(r), UserClaims(r), auth.UserGrant(req.Data.ID)) {
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	user := &data.VerifyUsers{
		UserID:               req.Data.ID,
		UserIDHash:           helpers.BytesToKeccak256Hash(common.HexToAddress(req.Data.ID).Bytes()),
		CreatedAt:            time.Now().UTC(),
		Status:               "not_verified",
		Proof:                []byte{},
		AgeLowerBound:        -1,
		ExpirationLowerBound: helpers.DefaultDateHex,
	}

	if req.Data.Attributes.Nationality != nil && *req.Data.Attributes.Nationality != "" {
		user.Nationality = *req.Data.Attributes.Nationality
	}

	if req.Data.Attributes.EventId != nil && *req.Data.Attributes.EventId != "" {
		user.EventID = *req.Data.Attributes.EventId
	}

	if req.Data.Attributes.AgeLowerBound != nil {
		user.AgeLowerBound = int(*req.Data.Attributes.AgeLowerBound)
	}

	if req.Data.Attributes.Uniqueness != nil {
		user.Uniqueness = *req.Data.Attributes.Uniqueness
	}

	if req.Data.Attributes.Sex != nil {
		user.SexEnable = *req.Data.Attributes.Sex
	}

	if req.Data.Attributes.ExpirationLowerBound != nil {
		user.ExpirationLowerBound = helpers.GetExpirationLowerBound(*req.Data.Attributes.ExpirationLowerBound)
	}

	dbUser, err := VerifyUsersQ(r).Upsert(user)
	if err != nil {
		Log(r).WithError(err).WithField("user", user).Errorf("failed to upsert user with userID [%s]", user.UserIDHash)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, responses.NewVerificationLinkLightResponse(dbUser, Callback(r).URL))

}
