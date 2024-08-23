package handlers

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	zkptypes "github.com/iden3/go-rapidsnark/types"
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

const (
	selector = "1"
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

	userIDHash, err := ExtractEventData(*proof)
	if err != nil {
		Log(r).WithError(err).Errorf("failed to get user hash from event data", userIDHash)
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	verifiedUser, err := VerifyUsersQ(r).WhereHashID(hex.EncodeToString(userIDHash[:])).Get()
	if err != nil {
		Log(r).WithError(err).Errorf("failed to get user with userHashID [%s]", userIDHash)
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if verifiedUser == nil {
		Log(r).Error("user is empty")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	selectorInt, err := strconv.Atoi(selector)
	if err != nil {
		fmt.Println("Error during conversion")
		return
	}
	if verifiedUser.Uniqueness {
		if proof.PubSignals[zk.TimestampUpperBound] == ProofParameters(r).TimestampUpperBound {
			if selectorInt&1<<9 == 0 {
				Log(r).Error("cannot extract timestampUpperBound from public signals")
				ape.RenderErr(w, problems.InternalError())
				return
			}
			if proof.PubSignals[zk.IdentityCounterUpperBound] == "0" {
				if selectorInt&1<<11 == 0 {
					Log(r).Error("cannot extract identityUpperBound from public signals")
					ape.RenderErr(w, problems.InternalError())
					return
				}
			}
		}
	}

	var verifyOpts = []zk.VerifyOption{
		zk.WithCitizenships(verifiedUser.Nationality),
		zk.WithProofSelectorValue(selector),
		zk.WithIdentitiesCounter(1),
		zk.WithAgeAbove(verifiedUser.AgeLowerBound),
		zk.WithEventID(ProofParameters(r).EventID),
	}
	err = Verifiers(r).Passport.VerifyProof(*proof, verifyOpts...)
	if err != nil {
		var vErr validation.Errors
		if !errors.As(err, &vErr) {
			verifiedUser.Status = "failed_verification"
			err = VerifyUsersQ(r).Update(verifiedUser)
			if err != nil {
				Log(r).WithError(err).Errorf("failed to update user status for userID [%s]", verifiedUser.UserIdHash)
				return
			}
			Log(r).WithError(err).Error("failed to verify proof")
			ape.RenderErr(w, problems.InternalError())
		}
		ape.RenderErr(w, problems.BadRequest(validation.Errors{"proof": err})...)
		return
	}
	verifiedUser.Status = "verified"
	err = VerifyUsersQ(r).Update(verifiedUser)
	if err != nil {
		Log(r).WithError(err).Errorf("failed to update user status for userID [%s]", verifiedUser.UserIdHash)
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
				Type: resources.RECEIVE_PROOF,
			},
			Attributes: resources.StatusAttributes{
				Status: user.Status,
			},
		},
	}
}

func ExtractEventData(proof zkptypes.ZKProof) ([]byte, error) {
	getter := zk.PubSignalGetter{Signals: proof.PubSignals, ProofType: zk.GlobalPassport}
	userIdHashDecimal, ok := new(big.Int).SetString(getter.Get(zk.EventData), 10)
	if !ok {
		return nil, fmt.Errorf("failed to parse event data")
	}
	var userIdHash [32]byte
	userIdHashDecimal.FillBytes(userIdHash[:])

	return userIdHash[:], nil
}
