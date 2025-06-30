package helpers

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

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
		timestampLowerBound = strconv.FormatInt(user.TimestampLowerBound.Time.Unix(), 10)
	}

	timestampUpperBound := "0"
	if user.TimestampUpperBound.Valid {
		timestampUpperBound = strconv.FormatInt(user.TimestampUpperBound.Time.Unix(), 10)
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
		Selector:                  strconv.FormatInt(int64(user.Selector), 10),
		TimestampLowerBound:       timestampLowerBound,
		TimestampUpperBound:       timestampUpperBound,
		CallbackUrl:               &callbackURL,
	}
}

func HexToDecimal(hexStr string) (string, error) {
	hexStr = strings.TrimPrefix(hexStr, "0x")

	decimal, ok := new(big.Int).SetString(hexStr, 16)
	if !ok {
		return "", fmt.Errorf("invalid hex format")
	}

	return decimal.String(), nil
}
