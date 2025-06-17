package helpers

import (
	"github.com/rarimo/verificator-svc/internal/data"
	"github.com/rarimo/verificator-svc/resources"
)

func BuildV2ProofParams(user *data.VerifyUsers, callbackURL string) resources.ProofParamsAttributes {
	birthDateLowerBound := DefaultDateHex
	if user.BirthDateLowerBound.Valid {
		birthDateLowerBound = user.BirthDateLowerBound.String
	}

	birthDateUpperBound := DefaultDateHex
	if user.BirthDateUpperBound.Valid {
		birthDateUpperBound = user.BirthDateUpperBound.String
	}

	eventData := user.UserIDHash
	if user.EventData.Valid {
		eventData = user.EventData.String
	}

	expirationDateUpperBound := DefaultDateHex
	if user.ExpirationDateUpperBound.Valid {
		expirationDateUpperBound = user.ExpirationDateUpperBound.String
	}

	timestampLowerBound := "0"
	if user.TimestampLowerBound.Valid {
		timestampLowerBound = user.TimestampLowerBound.String
	}

	timestampUpperBound := "0"
	if user.TimestampUpperBound.Valid {
		timestampUpperBound = user.TimestampUpperBound.String
	}

	return resources.ProofParamsAttributes{
		BirthDateLowerBound:       birthDateLowerBound,
		BirthDateUpperBound:       birthDateUpperBound,
		CitizenshipMask:           Utf8ToHex(user.Nationality),
		EventData:                 eventData,
		EventId:                   user.EventID,
		ExpirationDateLowerBound:  user.ExpirationLowerBound,
		ExpirationDateUpperBound:  expirationDateUpperBound,
		IdentityCounter:           user.IdentityCounter,
		IdentityCounterLowerBound: user.IdentityCounterLowerBound,
		IdentityCounterUpperBound: user.IdentityCounterUpperBound,
		Selector:                  user.Selector,
		TimestampLowerBound:       timestampLowerBound,
		TimestampUpperBound:       timestampUpperBound,
		CallbackUrl:               &callbackURL,
	}
}
