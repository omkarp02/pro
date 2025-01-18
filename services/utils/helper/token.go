package helper

import "github.com/omkarp02/pro/types"

var providerIdAndName = make(map[string]string)

func CreateAccessTokenPayload(id string, providerId string, role []string) types.ACCESS_TOKEN_PAYLOAD {
	return types.ACCESS_TOKEN_PAYLOAD{
		ID:         id,
		ProviderId: providerId,
		Role:       role,
	}
}

func CreateRefreshTokenPayload(id string, providerId string, role []string) types.REFRESH_TOKEN_PAYLOAD {
	return types.REFRESH_TOKEN_PAYLOAD{
		ID:         id,
		ProviderId: providerId,
		Role:       role,
	}
}

func ValidateDataForAccessToken(data interface{}) types.ACCESS_TOKEN_PAYLOAD {
	//here make function to check the data coming from local("user") is valid
	if data, ok := data.(types.ACCESS_TOKEN_PAYLOAD); ok {
		return data
	}

	panic("invalid data")
}
