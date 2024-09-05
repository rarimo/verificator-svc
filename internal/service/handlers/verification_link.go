package handlers

import (
	"github.com/rarimo/verificator-svc/internal/data"
	"github.com/rarimo/verificator-svc/internal/service/handlers/helpers"
	"github.com/rarimo/verificator-svc/internal/service/requests"
	"github.com/rarimo/verificator-svc/internal/service/responses"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
	"time"
)

func VerificationLink(w http.ResponseWriter, r *http.Request) {
	req, err := requests.VerificationLink(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	userIdHash, err := helpers.StringToPoseidonHash(req.Data.ID)
	if err != nil {
		Log(r).WithError(err).Errorf("failed to convert user with userID [%s] to poseidon hash", req.Data.ID)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	user := &data.VerifyUsers{
		UserID:        req.Data.ID,
		UserIDHash:    userIdHash,
		CreatedAt:     time.Now().UTC(),
		Status:        "not_verified",
		Proof:         []byte{},
		AgeLowerBound: -1,
	}

	if req.Data.Attributes.Nationality != nil && *req.Data.Attributes.Nationality != "" {
		user.Nationality = *req.Data.Attributes.Nationality
	}

	if req.Data.Attributes.EventId != nil && *req.Data.Attributes.EventId != "" {
		user.EventId = *req.Data.Attributes.EventId
	}

	if req.Data.Attributes.AgeLowerBound != nil {
		user.AgeLowerBound = int(*req.Data.Attributes.AgeLowerBound)
	}

	if req.Data.Attributes.Uniqueness != nil {
		user.Uniqueness = *req.Data.Attributes.Uniqueness
	}

	existingUser, err := VerifyUsersQ(r).WhereHashID(user.UserIDHash).Get()
	if err != nil {
		Log(r).WithError(err).Errorf("failed to query user with userID [%s]", userIdHash)
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if existingUser != nil {
		ape.Render(w, responses.NewVerificationLinkResponse(*existingUser, Callback(r).URL))
		return
	}

	if err = VerifyUsersQ(r).Insert(user); err != nil {
		Log(r).WithError(err).Errorf("failed to insert user with userID [%s]", user.UserIDHash)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, responses.NewVerificationLinkResponse(*user, Callback(r).URL))

}
