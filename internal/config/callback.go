package config

import (
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type CallbackConfiger interface {
	CallbackConfig() CallbackConfig
}

type CallbackConfig struct {
	Url string `fig:"url,required"`
}

type Callback struct {
	once   comfig.Once
	getter kv.Getter
}

func NewCallbackConfiger(getter kv.Getter) CallbackConfiger {
	return &Callback{
		getter: getter,
	}
}

func (p *Callback) CallbackConfig() CallbackConfig {
	return p.once.Do(func() interface{} {
		var cfg CallbackConfig
		err := figure.
			Out(&cfg).
			With(figure.BaseHooks).
			From(kv.MustGetStringMap(p.getter, "callback")).
			Please()

		if err != nil {
			panic(errors.Wrap(err, "failed to figure out callback"))
		}
		return cfg
	}).(CallbackConfig)
}
