package config

import (
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	zk "github.com/rarimo/zkverifier-kit"
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/kv"
)

const emptyETHAddr = "0x0000000000000000000000000000000000000000"

// 0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff
const maxEventId = 452312848583266388373324160190187140051835877600158453279131187530910662655

type Verifiers struct {
	Passport              *zk.Verifier
	ServiceStartTimestamp int64
	EventID               string
	Multiproof            bool
	PreserveUserIDCase    bool
	ERC1155               common.Address
}

func (c *config) Verifiers() Verifiers {
	return c.verifier.Do(func() interface{} {
		var cfg = struct {
			VerificationKeyPath      string         `fig:"verification_key_path,required"`
			AllowedIdentityTimestamp int64          `fig:"allowed_identity_timestamp,required"`
			EventID                  string         `fig:"event_id,required"`
			Multiproof               bool           `fig:"multiproof"`
			PreserveUserIDCase       bool           `fig:"preserve_user_id_case"`
			ERC1155                  common.Address `fig:"erc_1155"`
		}{
			ERC1155: common.HexToAddress(emptyETHAddr),
		}

		err := figure.
			Out(&cfg).
			With(figure.EthereumHooks).
			From(kv.MustGetStringMap(c.getter, "verifier")).
			Please()
		if err != nil {
			panic(fmt.Errorf("failed to figure out verifier: %w", err))
		}

		pass, err := zk.NewVerifier(nil,
			zk.WithProofType(zk.GlobalPassport),
			zk.WithVerificationKeyFile(cfg.VerificationKeyPath),
			zk.WithPassportRootVerifier(c.passport.ProvideVerifier()),
		)
		if err != nil {
			panic(fmt.Errorf("failed to initialize passport verifier: %w", err))
		}

		eventID, err := strconv.ParseInt(cfg.EventID, 10, 64)
		if err != nil {
			panic(fmt.Errorf("event_id must be valid int, %w", err))
		}

		if eventID > maxEventId {
			panic(fmt.Errorf("event_id must be less than 31 bytes"))
		}

		return Verifiers{
			Passport:              pass,
			ServiceStartTimestamp: cfg.AllowedIdentityTimestamp,
			EventID:               cfg.EventID,
			Multiproof:            cfg.Multiproof,
			PreserveUserIDCase:    cfg.PreserveUserIDCase,
			ERC1155:               cfg.ERC1155,
		}
	}).(Verifiers)
}
