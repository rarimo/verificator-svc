package config

import (
	"fmt"
	zk "github.com/rarimo/zkverifier-kit"
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/kv"
)

const (
	proofEventIDValue       = "proofEventIDValue"
	proofSelectorValue      = "proofSelectorValue"
	maxIdentityCount        = 1
	documentTypeID          = "ID"
	passportVerificationKey = "./proof_keys/passport.json"
)

type Verifiers struct {
	Passport *zk.Verifier
}

func (c *config) Verifiers() Verifiers {
	return c.verifier.Do(func() interface{} {
		var cfg struct {
			AllowedAge               int   `fig:"allowed_age,required"`
			AllowedIdentityTimestamp int64 `fig:"allowed_identity_timestamp,required"`
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
			zk.WithCitizenships("GLOBAL"),
			zk.WithVerificationKeyFile(passportVerificationKey),
			zk.WithAgeAbove(cfg.AllowedAge),
			zk.WithPassportRootVerifier(c.passport.ProvideVerifier()),
			zk.WithProofSelectorValue(proofSelectorValue),
			zk.WithEventID(proofEventIDValue),
			zk.WithIdentitiesCounter(maxIdentityCount),
			zk.WithIdentitiesCreationTimestampLimit(cfg.AllowedIdentityTimestamp),
			zk.WithDocumentType(documentTypeID),
		)
		if err != nil {
			panic(fmt.Errorf("failed to initialize passport verifier: %w", err))
		}

		return Verifiers{
			Passport: pass,
		}
	}).(Verifiers)
}
