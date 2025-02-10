package handlers

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/log"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/verificator-svc/internal/service/ctx"
	"github.com/rarimo/verificator-svc/internal/service/handlers/helpers"
	"github.com/rarimo/verificator-svc/internal/service/requests"
	"github.com/rarimo/verificator-svc/internal/service/responses"
	zk "github.com/rarimo/zkverifier-kit"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
)

func VerificationCallback(w http.ResponseWriter, r *http.Request) {
	req, err := requests.GetVerificationCallbackByID(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	var (
		proof   = req.Data.Attributes.Proof
		getter  = zk.PubSignalGetter{Signals: proof.PubSignals, ProofType: zk.GlobalPassport}
		eventID = ctx.Verifiers(r).EventID
	)

	proofJSON, err := json.Marshal(proof)
	if err != nil {
		ctx.Log(r).WithError(err).Errorf("failed to convert proof to json")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	identityCounterUpperBound, err := strconv.ParseInt(getter.Get(zk.IdentityCounterUpperBound), 10, 64)
	if err != nil {
		ctx.Log(r).WithError(err).Errorf("cannot extract identityUpperBound from public signals")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	identityTimestampUpperBound, err := strconv.ParseInt(getter.Get(zk.TimestampUpperBound), 10, 64)
	if err != nil {
		ctx.Log(r).WithError(err).Errorf("cannot extract identityUpperBound from public signals")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	selectorInt, err := strconv.Atoi(getter.Get(zk.Selector))
	if err != nil {
		ctx.Log(r).WithError(err).Errorf("cannot extract selector from public signals")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	userNationality, err := helpers.DecimalToHexToUtf8(getter.Get(zk.Citizenship))
	if err != nil {
		ctx.Log(r).WithError(err).Errorf("failed to convert decimal(nationality) to utf8")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	nullifier, ok := new(big.Int).SetString(getter.Get(zk.Nullifier), 10)
	if !ok {
		ctx.Log(r).Error("failed to parse nullifier")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	var nullifierBytes [32]byte
	nullifier.FillBytes(nullifierBytes[:])
	nullifierHex := hex.EncodeToString(nullifierBytes[:])

	userIDHash, err := helpers.ExtractEventData(getter)
	if err != nil {
		ctx.Log(r).WithError(err).Errorf("failed to extract user hash from event data")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"pub_signals/event_data": err,
		})...)
		return
	}

	verifiedUser, err := ctx.VerifyUsersQ(r).WhereHashID(userIDHash).Get()
	if err != nil {
		ctx.Log(r).WithError(err).Errorf("failed to get user with userHashID [%s]", userIDHash)
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	if verifiedUser == nil {
		ctx.Log(r).WithFields(logan.F{
			"event_data":   getter.Get(zk.EventData),
			"user_id_hash": userIDHash,
			"id":           req.Data.ID,
		}).Error("user not found or eventData != userHashID")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	verifiedUser.Status = "verified"

	if verifiedUser.EventID != "" {
		eventID = verifiedUser.EventID
	}

	if verifiedUser.Uniqueness {
		byNullifier, dbErr := ctx.VerifyUsersQ(r).FilterByNullifier(nullifierHex).Get()
		if dbErr != nil {
			ctx.Log(r).Error("Failed to get user by nullifier")
			ape.RenderErr(w, problems.BadRequest(dbErr)...)
			return
		}

		if !ctx.Verifiers(r).Multiproof && byNullifier != nil && byNullifier.UserIDHash != verifiedUser.UserIDHash {
			ctx.Log(r).WithError(err).Errorf("User with this nullifier [%s] but a different userIDHash already exists", nullifierHex)
			verifiedUser.Status = "uniqueness_check_failed"
		} else {
			verifiedUser.Nullifier = nullifierHex
		}
	}

	var verifyOpts = []zk.VerifyOption{
		zk.WithProofSelectorValue(getter.Get(zk.Selector)),
		zk.WithAgeAbove(verifiedUser.AgeLowerBound), // if not required -1
		zk.WithEventID(eventID),
	}
	if verifiedUser.Nationality != "" {
		verifyOpts = append(verifyOpts, zk.WithCitizenships(verifiedUser.Nationality))
	} else {
		verifiedUser.Nationality = userNationality
	}

	err = ctx.Verifiers(r).Passport.VerifyProof(proof, verifyOpts...)
	if err != nil {
		var vErr validation.Errors
		if errors.As(err, &vErr) {
			verifiedUser.Status = "failed_verification"
			ctx.Log(r).WithError(err).Error("failed to verify proof, updating user status")
			dbErr := ctx.VerifyUsersQ(r).Update(verifiedUser)
			if dbErr != nil {
				ctx.Log(r).WithError(dbErr).Errorf("failed to update user status for userID [%s]", verifiedUser.UserIDHash)
				ape.RenderErr(w, problems.InternalError())
				return
			}
			ape.RenderErr(w, problems.BadRequest(err)...)
			return
		}
		ctx.Log(r).WithError(err).Error("failed to verify proof")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if verifiedUser.Uniqueness {
		unique, uniquenessErr := helpers.CheckUniqueness(selectorInt, ctx.Verifiers(r).ServiceStartTimestamp, identityTimestampUpperBound, identityCounterUpperBound)
		if uniquenessErr != nil {
			ctx.Log(r).WithError(err).Errorf("failed to check uniqueness")
			ape.RenderErr(w, problems.BadRequest(uniquenessErr)...)
			return
		}
		if !unique {
			ctx.Log(r).WithFields(logan.F{
				"selector":                       selectorInt,
				"service_timestamp":              ctx.Verifiers(r).ServiceStartTimestamp,
				"identity_timestamp_upper_bound": identityTimestampUpperBound,
				"identity_counter_upper_bound":   identityCounterUpperBound,
				"user_id_hash":                   userIDHash,
			}).Errorf("failed to check uniqueness")
			verifiedUser.Status = "uniqueness_check_failed"
		}
	}

	verifiedUser.Proof = proofJSON
	err = ctx.VerifyUsersQ(r).Update(verifiedUser)
	if err != nil {
		ctx.Log(r).WithError(err).Errorf("failed to update user status for userID [%s]", verifiedUser.UserIDHash)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	log.Debug("Proof successfully verified")
	ape.Render(w, responses.NewVerificationCallbackResponse(*verifiedUser))
}
