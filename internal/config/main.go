package config

import (
	"github.com/rarimo/geo-auth-svc/pkg/auth"
	"github.com/rarimo/zkverifier-kit/root"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/copus"
	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type Config interface {
	comfig.Logger
	pgdb.Databaser
	types.Copuser
	comfig.Listenerer
	CallbackConfiger
	Verifiers() Verifiers
	SignatureVerificationConfiger
	auth.Auther //nolint:misspell
}

type config struct {
	comfig.Logger
	pgdb.Databaser
	types.Copuser
	comfig.Listenerer
	getter kv.Getter
	CallbackConfiger
	SignatureVerificationConfiger

	verifier    comfig.Once
	passport    root.VerifierProvider
	auth.Auther //nolint:misspell
}

func New(getter kv.Getter) Config {
	return &config{
		getter:                        getter,
		Databaser:                     pgdb.NewDatabaser(getter),
		Copuser:                       copus.NewCopuser(getter),
		Listenerer:                    comfig.NewListenerer(getter),
		Logger:                        comfig.NewLogger(getter, comfig.LoggerOpts{}),
		CallbackConfiger:              NewCallbackConfiger(getter),
		SignatureVerificationConfiger: NewSignatureVerificationConfiger(getter),
		passport:                      root.NewVerifierProvider(getter, root.PoseidonSMT),
		Auther:                        auth.NewAuther(getter), //nolint:misspell
	}
}
