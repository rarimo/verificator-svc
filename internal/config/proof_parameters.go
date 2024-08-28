package config

import (
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type ProofParametersConfiger interface {
	ProofParametersConfig() ProofParametersConfig
}

type ProofParametersConfig struct {
	EventID             string `fig:"event_id,required"`
	TimestampUpperBound string `fig:"timestamp_upper_bound,required"`
}

type ProofParameters struct {
	once   comfig.Once
	getter kv.Getter
}

func NewProofParametersConfiger(getter kv.Getter) ProofParametersConfiger {
	return &ProofParameters{
		getter: getter,
	}
}

func (p *ProofParameters) ProofParametersConfig() ProofParametersConfig {
	return p.once.Do(func() interface{} {
		var cfg ProofParametersConfig
		err := figure.
			Out(&cfg).
			With(figure.BaseHooks).
			From(kv.MustGetStringMap(p.getter, "proof_parameters")).
			Please()

		if err != nil {
			panic(errors.Wrap(err, "failed to figure out proof_parameters"))
		}
		return cfg
	}).(ProofParametersConfig)
}
