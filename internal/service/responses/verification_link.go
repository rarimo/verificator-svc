package responses

import (
	"fmt"
	"github.com/rarimo/verificator-svc/internal/data"
	"github.com/rarimo/verificator-svc/resources"
)

func NewVerificationLinkResponse(user data.VerifyUsers, host string) resources.LinksRequest {
	return resources.LinksRequest{
		Data: resources.Links{
			Key: resources.Key{
				ID:   user.UserID,
				Type: resources.VERIFICATION_LINK,
			},
			Attributes: resources.LinksAttributes{
				GetProofParams: fmt.Sprintf("%s/integrations/verificator-svc/public/proof-params/%s", host, user.UserIDHash),
			},
		},
	}
}
