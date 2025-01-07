package utils

import (
	"math/rand"
	"time"

	"github.com/omkarp02/pro/services/utils/helper"
	"github.com/omkarp02/pro/types"
	"github.com/omkarp02/pro/utils/constant"
	"github.com/omkarp02/pro/utils/errutil"
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func GetUserDataFromAccessClaimsData(claimsData interface{}) (types.ACCESS_TOKEN_PAYLOAD, error) {
	claimsMap, ok := claimsData.(map[string]interface{})
	if !ok {
		return types.ACCESS_TOKEN_PAYLOAD{}, errutil.InternalServerError("Invalid Format")
	}

	accessTokenPayload := helper.CreateAccessTokenPayload(claimsMap["ID"].(string), claimsMap["ProviderId"].(string))

	return accessTokenPayload, nil
}

func GetUserDataFromRefreshClaimsData(claimsData interface{}) (types.REFRESH_TOKEN_PAYLOAD, error) {
	claimsMap, ok := claimsData.(map[string]interface{})
	if !ok {
		return types.REFRESH_TOKEN_PAYLOAD{}, errutil.InternalServerError("Invalid Format")
	}

	refreshTokenPayload := helper.CreateRefreshTokenPayload(claimsMap["ID"].(string), claimsMap["ProviderId"].(string))

	return refreshTokenPayload, nil
}

func GenearteRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = constant.Charset[seededRand.Intn(len(constant.Charset))]
	}
	return string(b)
}

func GenerateRandomNumber(n int) int {
	res := 0
	temp := 0
	for i := 0; i < n; i++ {
		r := rand.Intn(10)
		res = r + temp
		temp = res * 10
	}
	return res
}
