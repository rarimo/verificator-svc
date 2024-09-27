package responses

import (
	"github.com/rarimo/verificator-svc/internal/data"
	"github.com/rarimo/verificator-svc/resources"
)

func NewGetUsersByIdResponse(user data.VerifyUsers) resources.UserParamsRequest {
	return resources.UserParamsRequest{
		Data: resources.UserParams{
			Key: resources.Key{
				ID:   user.UserID,
				Type: resources.USER,
			},
			Attributes: resources.UserParamsAttributes{
				AgeLowerBound: int32(user.AgeLowerBound),
				Nationality:   user.Nationality,
				Sex:           user.Sex,
			},
		},
	}
}
