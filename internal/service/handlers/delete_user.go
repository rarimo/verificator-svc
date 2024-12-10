package handlers

import (
	"net/http"

	"github.com/rarimo/verificator-svc/internal/service/requests"
	"github.com/rarimo/web3-auth-svc/pkg/auth"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	if !auth.Authenticates(UserClaims(r), auth.AdminGrant) {
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	userID, err := requests.DeleteUserByID(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	deletedUser, err := VerifyUsersQ(r).WhereID(userID).Get()
	if err != nil {
		Log(r).WithError(err).Error("failed to get user by userID")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if deletedUser == nil {
		Log(r).Debugf("User with userID=%s not found", userID)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if err = VerifyUsersQ(r).DeleteByID(deletedUser); err != nil {
		Log(r).Debugf("failed to delete user with UserID=%s", userID)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(204)
}
