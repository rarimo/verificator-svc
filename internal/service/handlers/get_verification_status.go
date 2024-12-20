package handlers

import (
	"net/http"

	"github.com/rarimo/verificator-svc/internal/service/ctx"
	"github.com/rarimo/verificator-svc/internal/service/requests"
	"github.com/rarimo/verificator-svc/internal/service/responses"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetVerificationStatusById(w http.ResponseWriter, r *http.Request) {
	userID, err := requests.GetPathUserID(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	verifiedUser, err := ctx.VerifyUsersQ(r).WhereID(userID).Get()
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get user by userID")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if verifiedUser == nil {
		ctx.Log(r).Debugf("User with userID=%s not found", userID)
		ape.RenderErr(w, problems.NotFound())
		return
	}

	ape.Render(w, responses.NewVerificationStatusByIdResponse(*verifiedUser))
}
