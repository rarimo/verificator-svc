package handlers

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/verificator-svc/internal/data"
	"github.com/rarimo/verificator-svc/internal/service/requests"
	"github.com/rarimo/verificator-svc/resources"
	zk "github.com/rarimo/zkverifier-kit"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"math/big"
	"net/http"
	"strconv"
)

func VerificationCallback(w http.ResponseWriter, r *http.Request) {
	req, err := requests.GetVerificationCallbackByID(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	var (
		proof  = req.Data.Attributes.Proof
		getter = zk.PubSignalGetter{Signals: proof.PubSignals, ProofType: zk.GlobalPassport}
	)

	proofJSON, err := json.Marshal(proof)
	if err != nil {
		Log(r).WithError(err).Errorf("failed to convert proof to json")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	userIDHash, err := ExtractEventData(getter)
	if err != nil {
		Log(r).WithError(err).Errorf("failed to extract user hash from event data")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"pub_signals/event_data": err,
		})...)
		return
	}

	verifiedUser, err := VerifyUsersQ(r).WhereHashID(userIDHash).Get()
	if err != nil {
		Log(r).WithError(err).Errorf("failed to get user with userHashID [%s]", userIDHash)
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if verifiedUser == nil {
		Log(r).Error("user not found or eventData != userHashID")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	selectorInt, err := strconv.Atoi(getter.Get(zk.Selector))
	if err != nil {
		Log(r).WithError(err).Errorf("cannot extract selector from public signals")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	identityCounterUpperBound, err := strconv.ParseInt(getter.Get(zk.IdentityCounterUpperBound), 10, 64)
	if err != nil {
		Log(r).WithError(err).Errorf("cannot extract identityUpperBound from public signals")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	// uniqueness check
	timestampUpperBoundMatches := getter.Get(zk.TimestampUpperBound) == ProofParameters(r).TimestampUpperBound
	timestampUpperBoundCheckRequired := selectorInt&(1<<timestampUpperBoundBit) == 0
	if verifiedUser.Uniqueness && timestampUpperBoundMatches {
		if timestampUpperBoundCheckRequired {
			Log(r).Error("uniqueness check failed: pub-signals.TimestampUpperBound != config.TimestampUpperBound and timestampUpperBoundBit == 0")
			ape.RenderErr(w, problems.InternalError())
			return
		}
	}
	identityCounterUpperBoundMatches := getter.Get(zk.IdentityCounterUpperBound) == "1"
	identityCounterUpperBoundCheckRequired := selectorInt&(1<<identityCounterUpperBoundBit) == 0
	if verifiedUser.Uniqueness && identityCounterUpperBoundMatches {
		if identityCounterUpperBoundCheckRequired {
			Log(r).Error("uniqueness check failed: pub-signals.IdentityCounterUpperBound != 1 and identityCounterUpperBoundBit == 0")
			ape.RenderErr(w, problems.InternalError())
			return
		}
	}

	var verifyOpts = []zk.VerifyOption{
		zk.WithCitizenships(verifiedUser.Nationality),
		zk.WithProofSelectorValue(getter.Get(zk.Selector)),
		zk.WithIdentitiesCounter(identityCounterUpperBound),
		zk.WithAgeAbove(verifiedUser.AgeLowerBound),
	}
	err = Verifiers(r).Passport.VerifyProof(*proof, verifyOpts...)
	if err != nil {
		var vErr validation.Errors
		if errors.As(err, &vErr) {
			verifiedUser.Status = "failed_verification"
			Log(r).WithError(err).Error("failed to verify proof, updating user status")
			dbErr := VerifyUsersQ(r).Update(verifiedUser)
			if dbErr != nil {
				Log(r).WithError(dbErr).Errorf("failed to update user status for userID [%s]", verifiedUser.UserIDHash)
				ape.RenderErr(w, problems.InternalError())
				return
			}
			ape.RenderErr(w, problems.BadRequest(err)...)
			return
		}
		Log(r).WithError(err).Error("failed to verify proof")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	verifiedUser.Status = "verified"
	verifiedUser.Proof = proofJSON
	err = VerifyUsersQ(r).Update(verifiedUser)
	if err != nil {
		Log(r).WithError(err).Errorf("failed to update user status for userID [%s]", verifiedUser.UserIDHash)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	log.Debug("Proof successfully verified")
	ape.Render(w, NewVerificationCallbackResponse(*verifiedUser))
}

func NewVerificationCallbackResponse(user data.VerifyUsers) resources.StatusResponse {
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

func ExtractEventData(getter zk.PubSignalGetter) (string, error) {
	userIDHashDecimal, ok := new(big.Int).SetString(getter.Get(zk.EventData), 10)
	if !ok {
		return "", fmt.Errorf("failed to parse event data")
	}
	var userIDHash [32]byte
	userIDHashDecimal.FillBytes(userIDHash[:])

	return fmt.Sprintf("0x%s", hex.EncodeToString(userIDHash[:])), nil
}
