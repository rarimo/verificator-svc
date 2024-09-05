package handlers

import (
	"github.com/rarimo/verificator-svc/internal/service/requests"
	"github.com/rarimo/verificator-svc/internal/service/responses"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

func GetVerificationStatusById(w http.ResponseWriter, r *http.Request) {
	userID, err := requests.GetVerificationStatusByID(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	verifiedUser, err := VerifyUsersQ(r).WhereID(userID).Get()
	if err != nil {
		Log(r).WithError(err).Error("failed to get user by userID")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if verifiedUser == nil {
		Log(r).Debugf("User with userID=%s not found", userID)
		ape.RenderErr(w, problems.NotFound())
		return
	}

	ape.Render(w, responses.NewVerificationStatusByIdResponse(*verifiedUser))
}
