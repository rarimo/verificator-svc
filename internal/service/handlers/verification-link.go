package handlers

import (
	"fmt"
	"github.com/rarimo/verificator-svc/internal/data"
	"github.com/rarimo/verificator-svc/internal/service/requests"
	"github.com/rarimo/verificator-svc/resources"
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

	userIdHash, err := StringToPoseidonHash(req.Data.ID)
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
		Nationality:   *req.Data.Attributes.Nationality,
		EventId:       *req.Data.Attributes.EventId,
		AgeLowerBound: int(*req.Data.Attributes.AgeLowerBound),
		Uniqueness:    *req.Data.Attributes.Uniqueness,
		Proof:         []byte{},
	}

	existingUser, err := VerifyUsersQ(r).WhereHashID(user.UserIDHash).Get()
	if err != nil {
		Log(r).WithError(err).Errorf("failed to query user with userID [%s]", userIdHash)
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if existingUser != nil {
		ape.Render(w, NewVerificationLinkResponse(*existingUser, Callback(r).URL))
		return
	}

	if err = VerifyUsersQ(r).Insert(user); err != nil {
		Log(r).WithError(err).Errorf("failed to insert user with userID [%s]", user.UserIDHash)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, NewVerificationLinkResponse(*user, Callback(r).URL))

}

func NewVerificationLinkResponse(user data.VerifyUsers, host string) resources.LinksRequest {
	return resources.LinksRequest{
		Data: resources.Links{
			Key: resources.Key{
				ID:   user.UserID,
				Type: resources.VERIFICATION_LINK,
			},
			Attributes: resources.LinksAttributes{
				CallbackUrl:    fmt.Sprintf("%s/integrations/verificator-svc/public/callback/%s", host, user.UserIDHash),
				GetProofParams: fmt.Sprintf("%s/integrations/verificator-svc/public/proof-params/%s", host, user.UserIDHash),
			},
		},
	}
}
