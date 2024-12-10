package helpers

import (
	"github.com/rarimo/web3-auth-svc/pkg/auth"
	"github.com/rarimo/web3-auth-svc/resources"
)

func Authenticates(client *auth.Client, claims []resources.Claim, grants ...auth.Grant) bool {
	if !client.Enabled {
		return true
	}

	return auth.Authenticates(claims, grants...)
}
