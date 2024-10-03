package handlers

import (
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum/log"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/verificator-svc/internal/service/handlers/helpers"
	"github.com/rarimo/verificator-svc/internal/service/requests"
	"github.com/rarimo/verificator-svc/internal/service/responses"
	zk "github.com/rarimo/zkverifier-kit"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
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
		proof   = req.Data.Attributes.Proof
		getter  = zk.PubSignalGetter{Signals: proof.PubSignals, ProofType: zk.GlobalPassport}
		eventID = Verifiers(r).EventID
	)

	proofJSON, err := json.Marshal(proof)
	if err != nil {
		Log(r).WithError(err).Errorf("failed to convert proof to json")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	identityCounterUpperBound, err := strconv.ParseInt(getter.Get(zk.IdentityCounterUpperBound), 10, 64)
	if err != nil {
		Log(r).WithError(err).Errorf("cannot extract identityUpperBound from public signals")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	identityTimestampUpperBound, err := strconv.ParseInt(getter.Get(zk.TimestampUpperBound), 10, 64)
	if err != nil {
		Log(r).WithError(err).Errorf("cannot extract identityUpperBound from public signals")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	selectorInt, err := strconv.Atoi(getter.Get(zk.Selector))
	if err != nil {
		Log(r).WithError(err).Errorf("cannot extract selector from public signals")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	userIDHash, err := helpers.ExtractEventData(getter)
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
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	if verifiedUser == nil {
		Log(r).Error("user not found or eventData != userHashID")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if verifiedUser.EventId != "" {
		eventID = verifiedUser.EventId
	}

	var verifyOpts = []zk.VerifyOption{
		zk.WithProofSelectorValue(getter.Get(zk.Selector)),
		zk.WithAgeAbove(verifiedUser.AgeLowerBound), // if not required -1
		zk.WithEventID(eventID),
	}
	if verifiedUser.Nationality != "" {
		verifyOpts = append(verifyOpts, zk.WithCitizenships(verifiedUser.Nationality))
	}

	err = Verifiers(r).Passport.VerifyProof(proof, verifyOpts...)
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
	verifiedUserNationality, err := helpers.DecimalToHexToUtf8(getter.Get(zk.Citizenship))
	if err != nil {
		Log(r).WithError(err).Errorf("failed to convert decimal(nationality) to utf8")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	verifiedUser.Nationality = verifiedUserNationality
	verifiedUser.Status = "verified"

	if verifiedUser.Uniqueness {
		status, uniquenessErr := helpers.CheckUniqueness(selectorInt, Verifiers(r).ServiceStartTimestamp, identityTimestampUpperBound, identityCounterUpperBound)
		if uniquenessErr != nil {
			Log(r).WithError(err).Errorf("failed to check uniqueness")
			ape.RenderErr(w, problems.BadRequest(uniquenessErr)...)
			return
		}
		verifiedUser.Status = status
	}

	verifiedUser.Proof = proofJSON
	err = VerifyUsersQ(r).Update(verifiedUser)
	if err != nil {
		Log(r).WithError(err).Errorf("failed to update user status for userID [%s]", verifiedUser.UserIDHash)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	log.Debug("Proof successfully verified")
	ape.Render(w, responses.NewVerificationCallbackResponse(*verifiedUser))
}
