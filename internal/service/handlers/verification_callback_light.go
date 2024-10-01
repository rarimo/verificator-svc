package handlers

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/ethereum/go-ethereum/log"
	"github.com/rarimo/verificator-svc/internal/service/handlers/helpers"
	"github.com/rarimo/verificator-svc/internal/service/requests"
	"github.com/rarimo/verificator-svc/internal/service/responses"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"math/big"
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

	validSignature := secp256k1.VerifySignature(pubKey, pubSignalsHash, signature[:64])
	if !validSignature {
		Log(r).Error("provided signature not valid")
		ape.RenderErr(w, problems.Unauthorized())
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

	userIDHashDecimal, ok := new(big.Int).SetString(pubSignals[10], 10)
	if !ok {
		Log(r).Error("failed to parse event data")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	var eventDataBytes [32]byte
	userIDHashDecimal.FillBytes(eventDataBytes[:])

	eventData := fmt.Sprintf("0x%s", hex.EncodeToString(eventDataBytes[:]))
	nationality, err := helpers.DecimalToHexToUtf8(pubSignals[6])
	if err != nil {
		Log(r).Error("failed to convert nationality from decimal to UTF8")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	sex, err := helpers.DecimalToHexToUtf8(pubSignals[7])
	if err != nil {
		Log(r).Error("failed to convert sex from decimal to UTF8")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if verifiedUser.Nationality == "" && pubSignals[6] != "0" {
		verifiedUser.Nationality = nationality
	}
	if verifiedUser.Sex == "" && pubSignals[7] != "0" {
		verifiedUser.Sex = sex
	}

	verifiedUser.Status = "verified"
	if eventData != userIDHash {
		Log(r).WithError(err).Errorf("failed to verify user: EventData from pub-signals [%s] != userIdHash from db [%s]", eventData, userIDHash)
		verifiedUser.Status = "failed_verification"
	}
	if verifiedUser.Nationality != nationality {
		Log(r).WithError(err).Errorf("failed to verify user with UserIdHash[%s]: Citizenship from pub-signals [%s] != User.Citizenship from db [%s]", userIDHash, nationality, verifiedUser.Nationality)
		verifiedUser.Status = "failed_verification"
	}
	if verifiedUser.Sex != sex {
		Log(r).WithError(err).Errorf("failed to verify user with UserIdHash[%s]: Sex from pub-signals [%s] != User.Sex from db [%s]", userIDHash, sex, verifiedUser.Sex)
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
