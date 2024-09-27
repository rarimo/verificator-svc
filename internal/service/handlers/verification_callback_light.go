package handlers

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/ethereum/go-ethereum/log"
	"github.com/rarimo/verificator-svc/internal/service/handlers/helpers"
	"github.com/rarimo/verificator-svc/internal/service/requests"
	"github.com/rarimo/verificator-svc/internal/service/responses"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

func VerificationSignatureCallback(w http.ResponseWriter, r *http.Request) {
	req, err := requests.GetVerificationCallbackSignatureByID(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	var (
		userIDHash = req.Data.ID
		pubSignals = req.Data.Attributes.PubSignals
	)

	signature, err := hex.DecodeString(req.Data.Attributes.Signature)
	if err != nil {
		Log(r).Error("cannot decode signature from string to bytes")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	pubKey, err := hex.DecodeString(SignatureVerification(r).PubKey)
	if err != nil {
		Log(r).Error("cannot decode public-key from string to bytes")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	pubSignalsHash, err := helpers.PubSignalsToSha256(pubSignals)
	if err != nil {
		Log(r).Error("failed to convert pubSignal array to sha256")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	verifiedUser, err := VerifyUsersQ(r).WhereHashID(userIDHash).Get()
	if err != nil {
		Log(r).WithError(err).Errorf("failed to get user with userHashID [%s]", userIDHash)
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	if verifiedUser == nil {
		Log(r).Error("user not found or eventData != userHashID")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	//if pubSignals[10] != userIDHash {
	//	Log(r).Error("failed to verify eventData")
	//	ape.RenderErr(w, problems.NotFound())
	//	return
	//}
	//if pubSignals[22] == "0" && pubSignals[22] != verifiedUser.Nationality {
	//	Log(r).Error("failed to verify citizenship")
	//	ape.RenderErr(w, problems.NotFound())
	//	return
	//}
	//if verifiedUser.SexEnable && pubSignals[7] == "0" {
	//	Log(r).Error("failed to verify sex")
	//	ape.RenderErr(w, problems.NotFound())
	//	return
	//}

	verificationStatus := secp256k1.VerifySignature(pubKey, pubSignalsHash, signature[:64])
	if verificationStatus {
		verifiedUser.Status = "verified"
	} else {
		verifiedUser.Status = "failed_verification"
	}

	err = VerifyUsersQ(r).Update(verifiedUser)
	if err != nil {
		Log(r).WithError(err).Errorf("failed to update user status for userID [%s]", verifiedUser.UserIDHash)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	log.Debug("Signature successfully verified")
	ape.Render(w, responses.NewVerificationCallbackResponse(*verifiedUser))
}
