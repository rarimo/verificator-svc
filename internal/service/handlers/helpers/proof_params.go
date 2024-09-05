package helpers

import (
	"encoding/hex"
	"fmt"
	"github.com/iden3/go-iden3-crypto/poseidon"
	zk "github.com/rarimo/zkverifier-kit"
	"github.com/status-im/keycard-go/hexutils"
	"math/big"
	"time"
)

const (
	NullifierBit                 = 0
	CitizenshipBit               = 5
	TimestampUpperBoundBit       = 9
	IdentityCounterUpperBoundBit = 11
	ExpirationDateLowerboundBit  = 12
	ExpirationDateUpperbound     = 13
	BirthDateLowerboundBit       = 14
	BirthDateUpperboundBit       = 15
	BirthDateFormat              = "060102"
)

func StringToPoseidonHash(inputString string) (string, error) {
	inputBytes := []byte(inputString)

	hash, err := poseidon.HashBytes(inputBytes)
	if err != nil {
		return "", fmt.Errorf("failde to convert input bytes to hash: %w", err)

	}
	return fmt.Sprintf("0x%s", hex.EncodeToString(hash.Bytes())), nil
}

func Utf8ToHex(input string) string {
	bytes := []byte(input)
	hexString := hexutils.BytesToHex(bytes)
	return fmt.Sprintf("0x%s", hexString)
}

func CalculateBirthDateHex(ageLowerBound int) string {
	allowedBirthDate := time.Now().UTC().AddDate(-ageLowerBound, 0, 0)
	formattedDate := []byte(allowedBirthDate.Format(BirthDateFormat))
	hexBirthDateLoweBound := hexutils.BytesToHex(formattedDate)

	return fmt.Sprintf("0x%s", hexBirthDateLoweBound)
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

func CalculateProofSelector(uniqueness bool, ageLowerBound int, nationality string) int {
	var bitLine uint32
	bitLine |= 1 << NullifierBit
	if nationality != "" {
		bitLine |= 1 << CitizenshipBit
	}
	if ageLowerBound != -1 {
		bitLine |= 1 << BirthDateUpperboundBit
	}
	if uniqueness {
		bitLine |= 1 << TimestampUpperBoundBit
		bitLine |= 1 << IdentityCounterUpperBoundBit
	}

	return int(bitLine)
}

func CheckUniqueness(selectorInt int, serviceStartTimestamp, identityTimestampUpperBound, identityCounterUpperBound int64) (string, error) {
	var (
		status           = "uniqueness_check_failed"
		timestampSuccess = false
		counterSuccess   = false
	)

	if selectorInt&(1<<TimestampUpperBoundBit) == 0 && selectorInt&(1<<IdentityCounterUpperBoundBit) == 0 {
		return "", fmt.Errorf("both timestampUpperBoundBit and identityCounterUpperBoundBit are not set in selector")
	}
	if selectorInt&(1<<TimestampUpperBoundBit) == 1<<TimestampUpperBoundBit {
		if identityTimestampUpperBound <= serviceStartTimestamp {
			timestampSuccess = true
		}
	}
	if selectorInt&(1<<IdentityCounterUpperBoundBit) == 1<<IdentityCounterUpperBoundBit {
		if identityCounterUpperBound <= 1 {
			counterSuccess = true
		}
	}
	if timestampSuccess || counterSuccess {
		return "verified", nil
	}

	return status, nil
}
