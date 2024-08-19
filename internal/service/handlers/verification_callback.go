package handlers

import (
	"errors"
	"github.com/ethereum/go-ethereum/log"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/verificator-svc/internal/service/requests"
	zk "github.com/rarimo/zkverifier-kit"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

func VerificationCallback(w http.ResponseWriter, r *http.Request) {
	req, err := requests.GetVerificationCallbackByID(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	var (
		proof = req.Data.Attributes.Proof
	)

	if proof == nil {
		log.Debug("Proof is not provided")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if proof.PubSignals == nil || len(proof.PubSignals) == 0 || len(proof.PubSignals) != 22 {
		log.Debug("PubSignals is not provided or empty")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	getter := zk.PubSignalGetter{Signals: proof.PubSignals, ProofType: zk.GlobalPassport}
	userIdHash := getter.Get(zk.EventData)
	err = Verifiers(r).Passport.VerifyProof(*proof)
	if err != nil {
		var vErr validation.Errors
		if !errors.As(err, &vErr) {
			Log(r).WithError(err).Error("Failed to verify proof")
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}
	verifiedUser, err := VerifyUsersQ(r).WhereHashID(userIdHash).Get()
	if err != nil {
		Log(r).WithError(err).Error("Failed to get user by userHashId")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	verifiedUser.Status = "true"
	err = VerifyUsersQ(r).Update(verifiedUser)
	if err != nil {
		Log(r).WithError(err).Errorf("failed to update user status for userID [%s]", verifiedUser.UserIdHash)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	log.Debug("Proof successfully verified")
	ape.Render(w, "Success")
}
