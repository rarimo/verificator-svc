package handlers

import (
	"github.com/rarimo/verificator-svc/internal/data"
	"github.com/rarimo/verificator-svc/internal/service/requests"
	"github.com/rarimo/verificator-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

func GetVerificationStatusById(w http.ResponseWriter, r *http.Request) {
	userId, err := requests.GetVerificationStatusByID(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	verifiedUser, err := VerifyUsersQ(r).WhereID(userId).Get()
	if err != nil {
		Log(r).WithError(err).Error("Failed to get user by userId")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if verifiedUser == nil {
		Log(r).Debugf("User for userId=%s not found", userId)
		ape.RenderErr(w, problems.NotFound())
		return
	}

	ape.Render(w, NewVerificationStatusByIdResponse(*verifiedUser))
}

func NewVerificationStatusByIdResponse(user data.VerifyUsers) resources.StatusResponse {
	return resources.StatusResponse{
		Data: resources.Status{
			Key: resources.Key{
				ID:   user.UserID,
				Type: resources.USER_STATUS,
			},
			Attributes: resources.StatusAttributes{
				Status: user.Status,
			},
		},
	}
}
