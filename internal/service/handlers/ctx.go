package handlers

import (
	"context"
	"github.com/rarimo/verificator-svc/internal/config"
	"github.com/rarimo/verificator-svc/internal/data"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	verifyUserQCtxKey
	verifiersCtxKey
	callbackCtxKey
	proofParametersCtxKey
	signatureVerificationCtxKey
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func CtxVerifyUsersQ(q data.VerifyUsersQ) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, verifyUserQCtxKey, q)
	}
}

func VerifyUsersQ(r *http.Request) data.VerifyUsersQ {
	return r.Context().Value(verifyUserQCtxKey).(data.VerifyUsersQ).New()
}

func CtxVerifiers(v config.Verifiers) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, verifiersCtxKey, v)
	}
}

func Verifiers(r *http.Request) config.Verifiers {
	return r.Context().Value(verifiersCtxKey).(config.Verifiers)
}

func CtxCallback(c config.CallbackConfig) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, callbackCtxKey, c)
	}
}

func Callback(r *http.Request) config.CallbackConfig {
	return r.Context().Value(callbackCtxKey).(config.CallbackConfig)
}

func CtxSignatureVerification(c config.SignatureVerificationConfig) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, signatureVerificationCtxKey, c)
	}
}

func SignatureVerification(r *http.Request) config.SignatureVerificationConfig {
	return r.Context().Value(signatureVerificationCtxKey).(config.SignatureVerificationConfig)
}
