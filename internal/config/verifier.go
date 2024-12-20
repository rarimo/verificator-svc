package config

import (
	"fmt"

	zk "github.com/rarimo/zkverifier-kit"
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/kv"
)

type Verifiers struct {
	Passport              *zk.Verifier
	ServiceStartTimestamp int64
	EventID               string
	Multiproof            bool
	PreserveUserIDCase    bool
}

func (c *config) Verifiers() Verifiers {
	return c.verifier.Do(func() interface{} {
		var cfg struct {
			VerificationKeyPath      string `fig:"verification_key_path,required"`
			AllowedIdentityTimestamp int64  `fig:"allowed_identity_timestamp,required"`
			EventID                  string `fig:"event_id,required"`
			Multiproof               bool   `fig:"multiproof"`
			PreserveUserIDCase       bool   `fig:"preserve_user_id_case"`
		}

		err := figure.
			Out(&cfg).
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

		return Verifiers{
			Passport:              pass,
			ServiceStartTimestamp: cfg.AllowedIdentityTimestamp,
			EventID:               cfg.EventID,
			Multiproof:            cfg.Multiproof,
			PreserveUserIDCase:    cfg.PreserveUserIDCase,
		}
	}).(Verifiers)
}
