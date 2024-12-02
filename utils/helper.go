package utils

import (
	"errors"

	"github.com/omkarp02/pro/types"
)

func GetUserDataFromAccessClaimsData(claimsData interface{}) (types.ACCESS_TOKEN_PAYLOAD, error) {
	claimsMap, ok := claimsData.(map[string]interface{})
	if !ok {
		return types.ACCESS_TOKEN_PAYLOAD{}, errors.New("Invalid Format")
	}

	return types.ACCESS_TOKEN_PAYLOAD{
		ID: claimsMap["ID"].(string),
	}, nil
}

func GetUserDataFromRefreshClaimsData(claimsData interface{}) (types.REFRESH_TOKEN_PAYLOAD, error) {
	claimsMap, ok := claimsData.(map[string]interface{})
	if !ok {
		return types.REFRESH_TOKEN_PAYLOAD{}, errors.New("Invalid Format")
	}

	return types.REFRESH_TOKEN_PAYLOAD{
		ID: claimsMap["ID"].(string),
	}, nil
}
