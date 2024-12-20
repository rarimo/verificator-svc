package handlers

import (
	"net/http"

	"github.com/rarimo/verificator-svc/internal/service/ctx"
	"github.com/rarimo/verificator-svc/internal/service/handlers/helpers"
	"github.com/rarimo/verificator-svc/internal/service/requests"
	"github.com/rarimo/web3-auth-svc/pkg/auth"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	if !auth.Authenticates(ctx.UserClaims(r), auth.AdminGrant) {
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	userID, err := requests.GetPathUserID(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if !helpers.Authenticates(ctx.AuthClient(r), ctx.UserClaims(r), auth.UserGrant(userID)) {
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	deletedUser, err := ctx.VerifyUsersQ(r).WhereID(userID).Get()
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get user by userID")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if deletedUser == nil {
		ctx.Log(r).Debugf("User with userID=%s not found", userID)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if err = ctx.VerifyUsersQ(r).DeleteByID(deletedUser); err != nil {
		ctx.Log(r).Debugf("failed to delete user with UserID=%s", userID)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(204)
}
