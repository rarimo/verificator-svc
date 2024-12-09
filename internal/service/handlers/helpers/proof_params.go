package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/iden3/go-iden3-crypto/poseidon"
	zk "github.com/rarimo/zkverifier-kit"
	"github.com/status-im/keycard-go/hexutils"
)

const (
	NullifierBit                 = 0
	CitizenshipBit               = 5
	SexBit                       = 6
	TimestampUpperBoundBit       = 9
	IdentityCounterUpperBoundBit = 11
	ExpirationDateLowerboundBit  = 12
	ExpirationDateUpperbound     = 13
	BirthDateLowerboundBit       = 14
	BirthDateUpperboundBit       = 15
	DateFormat                   = "060102"
)

type SelectorParams struct {
	Uniqueness           bool
	AgeLowerBound        int
	Nationality          string
	SexEnable            bool
	NationalityEnable    bool
	ExpirationLowerBound bool
}

func PubSignalsToSha256(pubSignals []string) ([]byte, error) {
	var hash = sha256.New()
	for _, pubSignalByte := range pubSignals {
		if len(pubSignalByte) > 1 && pubSignalByte[:2] == "0x" {
			pubSignalBytes, convertErr := hex.DecodeString(pubSignalByte[2:])
			if convertErr != nil {
				return nil, fmt.Errorf("error in converting pubSignalHex: %v", pubSignalByte)
			}
			hash.Write(pubSignalBytes)
		} else {
			pubSignalDecimal, ok := new(big.Int).SetString(pubSignalByte, 10)
			if !ok {
				return nil, fmt.Errorf("error in converting pubSignal: %v", pubSignalByte)
			}
			hash.Write(pubSignalDecimal.Bytes())
		}
	}
	messageHash := hash.Sum(nil)

	return messageHash, nil
}

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

func DecimalToHexToUtf8(input string) (string, error) {
	inputBig, ok := new(big.Int).SetString(input, 10)
	if !ok {
		return "", fmt.Errorf("failed to parse big int when converting to UTF8")
	}

	inputUtf8 := string(inputBig.Bytes())

	return inputUtf8, nil
}

func CalculateBirthDateHex(ageLowerBound int) string {
	return FormatDateTime(time.Now().UTC().AddDate(-ageLowerBound, 0, 0))
}

func GetExpirationLowerBound(expirationLowerBound bool) string {
	if !expirationLowerBound {
		return "52983525027888"
	}

	return FormatDateTime(time.Now().UTC())
}

func FormatDateTime(date time.Time) string {
	return fmt.Sprintf("0x%s", hexutils.BytesToHex([]byte(date.Format(DateFormat))))
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

func CalculateProofSelector(p SelectorParams) int {
	var bitLine uint32
	bitLine |= 1 << NullifierBit
	if p.Nationality != "" || p.NationalityEnable {
		bitLine |= 1 << CitizenshipBit
	}
	if p.SexEnable {
		bitLine |= 1 << SexBit
	}
	if p.AgeLowerBound != -1 {
		bitLine |= 1 << BirthDateUpperboundBit
	}
	if p.Uniqueness {
		bitLine |= 1 << TimestampUpperBoundBit
		bitLine |= 1 << IdentityCounterUpperBoundBit
	}
	if p.ExpirationLowerBound {
		bitLine |= 1 << ExpirationDateLowerboundBit
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
