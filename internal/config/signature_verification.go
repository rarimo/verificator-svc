package config

import (
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type SignatureVerificationConfiger interface {
	SignatureVerificationConfig() SignatureVerificationConfig
}

type SignatureVerificationConfig struct {
	PubKey string `fig:"pub_key,required"`
}

type SignatureVerification struct {
	once   comfig.Once
	getter kv.Getter
}

func NewSignatureVerificationConfiger(getter kv.Getter) SignatureVerificationConfiger {
	return &SignatureVerification{
		getter: getter,
	}
}

func (p *SignatureVerification) SignatureVerificationConfig() SignatureVerificationConfig {
	return p.once.Do(func() interface{} {
		var cfg SignatureVerificationConfig
		err := figure.
			Out(&cfg).
			With(figure.BaseHooks).
			From(kv.MustGetStringMap(p.getter, "signature_verification")).
			Please()

		if err != nil {
			panic(errors.Wrap(err, "failed to figure out callback"))
		}
		return cfg
	}).(SignatureVerificationConfig)
}
