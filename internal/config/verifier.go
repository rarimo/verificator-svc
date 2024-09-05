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
}

func (c *config) Verifiers() Verifiers {
	return c.verifier.Do(func() interface{} {
		var cfg struct {
			VerificationKeyPath      string `fig:"verification_key_path,required"`
			AllowedIdentityTimestamp int64  `fig:"allowed_identity_timestamp,required"`
			EventID                  string `fig:"event_id,required"`
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
		}
	}).(Verifiers)
}
