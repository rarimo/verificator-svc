package helpers

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestHash(t *testing.T) {
	hash := BytesToKeccak256Hash(common.HexToAddress("0x52749da41B7196A7001D85Ce38fa794FE0F9044E").Bytes())
	fmt.Println("HASH", hash)
	hashBytes, _ := hex.DecodeString(hash[2:])

	fmt.Println("HASH BYTES LEN", len(hashBytes))

	var eventDataBytes [32]byte
	copy(eventDataBytes[:], hashBytes)

	fmt.Printf("0x%s\n", hex.EncodeToString(eventDataBytes[:]))
}

func TestSelector(t *testing.T) {
	fmt.Println(CheckUniqueness(2561, 1726672309, 1727256960, 2))
}
