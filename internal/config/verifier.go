package config

import (
	"fmt"

	zk "github.com/rarimo/zkverifier-kit"
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/kv"
)

const (
	passportVerificationKey = "./proof_keys/passport.json"
)

type Verifiers struct {
	Passport *zk.Verifier
}

func (c *config) Verifiers() Verifiers {
	return c.verifier.Do(func() interface{} {
		var cfg struct {
			AllowedAge               int    `fig:"allowed_age,required"`
			VerificationKeyPath      string `fig:"verification_key_path,required"`
			AllowedIdentityTimestamp int64  `fig:"allowed_identity_timestamp,required"`
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
			zk.WithIdentitiesCreationTimestampLimit(cfg.AllowedIdentityTimestamp),
		)
		if err != nil {
			panic(fmt.Errorf("failed to initialize passport verifier: %w", err))
		}

		return Verifiers{
			Passport: pass,
		}
	}).(Verifiers)
}
