package handlers

import (
	"context"
	"github.com/rarimo/verificator-svc/internal/data"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	verifyUserQCtxKey
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func VerifyUsersQ(r *http.Request) data.VerifyUsersQ {
	return r.Context().Value(verifyUserQCtxKey).(data.VerifyUsersQ).New()
}
