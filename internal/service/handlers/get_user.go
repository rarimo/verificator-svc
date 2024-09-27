package handlers

import (
	"github.com/rarimo/verificator-svc/internal/service/requests"
	"github.com/rarimo/verificator-svc/internal/service/responses"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	userID, err := requests.DeleteUserByID(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	user, err := VerifyUsersQ(r).WhereID(userID).Get()
	if err != nil {
		Log(r).WithError(err).Error("failed to get user by userID")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if user == nil {
		Log(r).Debugf("User with userID=%s not found", userID)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, responses.NewGetUsersByIdResponse(*user))
}
