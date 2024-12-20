package handlers

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/ethereum/go-ethereum/log"
	"github.com/rarimo/verificator-svc/internal/service/ctx"
	"github.com/rarimo/verificator-svc/internal/service/handlers/helpers"
	"github.com/rarimo/verificator-svc/internal/service/requests"
	"github.com/rarimo/verificator-svc/internal/service/responses"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
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
		ctx.Log(r).Error("cannot decode signature from string to bytes")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	pubKey, err := hex.DecodeString(ctx.SignatureVerification(r).PubKey)
	if err != nil {
		ctx.Log(r).Error("cannot decode public-key from string to bytes")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	pubSignalsHash, err := helpers.PubSignalsToSha256(pubSignals)
	if err != nil {
		ctx.Log(r).Error("failed to convert pubSignal array to sha256")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	validSignature := secp256k1.VerifySignature(pubKey, pubSignalsHash, signature[:64])
	if !validSignature {
		ctx.Log(r).Error("provided signature not valid")
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	userIDHashDecimal, ok := new(big.Int).SetString(pubSignals[10], 10)
	if !ok {
		ctx.Log(r).Error("failed to parse event data")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	var eventDataBytes [31]byte
	userIDHashDecimal.FillBytes(eventDataBytes[:])

	eventData := fmt.Sprintf("0x%s", userIDHashDecimal.Text(16))
	nationality, err := helpers.DecimalToHexToUtf8(pubSignals[6])
	if err != nil {
		ctx.Log(r).Error("failed to convert nationality from decimal to UTF8")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	sex, err := helpers.DecimalToHexToUtf8(pubSignals[7])
	if err != nil {
		ctx.Log(r).Error("failed to convert sex from decimal to UTF8")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	sex = transformGender(sex)

	nullifier, ok := new(big.Int).SetString(pubSignals[0], 10)
	if !ok {
		ctx.Log(r).Error("failed to parse nullifier")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	var nullifierBytes [32]byte
	nullifier.FillBytes(nullifierBytes[:])
	nullifierHex := hex.EncodeToString(nullifierBytes[:])

	anonymousID, ok := new(big.Int).SetString(pubSignals[11], 10)
	if !ok {
		ctx.Log(r).Error("failed to parse anonymous_id")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	var anonymousIDBytes [32]byte
	anonymousID.FillBytes(anonymousIDBytes[:])
	anonymousIDHex := hex.EncodeToString(anonymousIDBytes[:])

	verifiedUser, err := ctx.VerifyUsersQ(r).WhereHashID(userIDHash).Get()
	if err != nil {
		ctx.Log(r).WithError(err).Errorf("failed to get user with userHashID [%s]", userIDHash)
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	if verifiedUser == nil {
		ctx.Log(r).Error("user not found or eventData != userHashID")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if verifiedUser.Nationality == "" && pubSignals[6] != "0" {
		verifiedUser.Nationality = nationality
	}
	if verifiedUser.Sex == "" && pubSignals[7] != "0" {
		verifiedUser.Sex = sex
	}

	verifiedUser.Status = "verified"

	if verifiedUser.Uniqueness {
		byAnonymousID, dbErr := ctx.VerifyUsersQ(r).FilterByInternalAID(anonymousIDHex).Get()
		if dbErr != nil {
			ctx.Log(r).Error("Failed to get user by anonymous_id")
			ape.RenderErr(w, problems.BadRequest(dbErr)...)
			return
		}

		byNullifier, dbErr := ctx.VerifyUsersQ(r).FilterByNullifier(nullifierHex).Get()
		if dbErr != nil {
			ctx.Log(r).Error("Failed to get user by nullifier")
			ape.RenderErr(w, problems.BadRequest(dbErr)...)
			return
		}

		if !ctx.Verifiers(r).Multiproof && byAnonymousID != nil {
			if byAnonymousID.UserIDHash != verifiedUser.UserIDHash {
				ctx.Log(r).WithError(err).Errorf("User with anonymous_id [%s] but a different userIDHash already exists", anonymousIDHex)
				verifiedUser.Status = "failed_verification"
			}
		} else {
			verifiedUser.AnonymousID = anonymousIDHex
		}

		if !ctx.Verifiers(r).Multiproof && byNullifier != nil {
			if byNullifier.UserIDHash != verifiedUser.UserIDHash {
				ctx.Log(r).WithError(err).Errorf("User with nullifier [%s] but a different userIDHash already exists", nullifierHex)
				verifiedUser.Status = "failed_verification"
			}
		} else {
			verifiedUser.Nullifier = nullifierHex
		}
	}

	if eventData != userIDHash {
		ctx.Log(r).WithError(err).Errorf("failed to verify user: EventData from pub-signals [%s] != userIdHash from db [%s]", eventData, userIDHash)
		verifiedUser.Status = "failed_verification"
	}
	if verifiedUser.Nationality != nationality {
		ctx.Log(r).WithError(err).Errorf("failed to verify user with UserIdHash[%s]: Citizenship from pub-signals [%s] != User.Citizenship from db [%s]", userIDHash, nationality, verifiedUser.Nationality)
		verifiedUser.Status = "failed_verification"
	}
	if verifiedUser.Sex != sex {
		ctx.Log(r).WithError(err).Errorf("failed to verify user with UserIdHash[%s]: Sex from pub-signals [%s] != User.Sex from db [%s]", userIDHash, sex, verifiedUser.Sex)
		verifiedUser.Status = "failed_verification"
	}

	err = ctx.VerifyUsersQ(r).Update(verifiedUser)
	if err != nil {
		ctx.Log(r).WithError(err).Errorf("failed to update user status for userID [%s]", verifiedUser.UserIDHash)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	log.Debug("Signature successfully verified")
	ape.Render(w, responses.NewVerificationCallbackResponse(*verifiedUser))
}

func transformGender(input string) string {
	switch strings.ToUpper(input) {
	case "MALE":
		return "M"
	case "FEMALE":
		return "F"
	case "OTHERS":
		return "O"
	// assume that other cases are valid (M,F,O)
	default:
		return input
	}
}
