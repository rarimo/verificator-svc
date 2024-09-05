package responses

import (
	"github.com/rarimo/verificator-svc/internal/data"
	"github.com/rarimo/verificator-svc/resources"
)

func NewVerificationStatusByIdResponse(user data.VerifyUsers) resources.StatusResponse {
	return resources.StatusResponse{
		Data: resources.Status{
			Key: resources.Key{
				ID:   user.UserID,
				Type: resources.USER_STATUS,
			},
			Attributes: resources.StatusAttributes{
				Status: user.Status,
			},
		},
	}
}
