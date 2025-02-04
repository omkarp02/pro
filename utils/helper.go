package utils

import (
	cryptorand "crypto/rand"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"strconv"
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

	accessTokenPayload := helper.CreateAccessTokenPayload(formatClaimsData(claimsMap))

	fmt.Println(accessTokenPayload, "<<<<<<<<< here is the token")

	return accessTokenPayload, nil
}

func GetUserDataFromRefreshClaimsData(claimsData interface{}) (types.REFRESH_TOKEN_PAYLOAD, error) {
	claimsMap, ok := claimsData.(map[string]interface{})
	if !ok {
		return types.REFRESH_TOKEN_PAYLOAD{}, errutil.InternalServerError("Invalid Format")
	}

	refreshTokenPayload := helper.CreateRefreshTokenPayload(formatClaimsData(claimsMap))

	return refreshTokenPayload, nil
}

func formatClaimsData(claimsMap map[string]interface{}) (string, string, []string) {

	roles := claimsMap["Role"].([]interface{})
	formattedRoles := []string{}

	for _, v := range roles {
		if str, ok := v.(string); ok { // Type assertion to check if the element is a string
			formattedRoles = append(formattedRoles, str)
		}
	}
	return claimsMap["ID"].(string), claimsMap["ProviderId"].(string), formattedRoles
}

func GenearteRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = constant.Charset[seededRand.Intn(len(constant.Charset))]
	}
	return string(b)
}

func GenerateRandomNumber(n int) int {
	// res := 0
	// temp := 0
	// for i := 0; i < n; i++ {
	// 	r := rand.Intn(10)
	// 	res = r + temp
	// 	temp = res * 10
	// }
	// return res
	number := int64(math.Pow(10, float64(n-1)))

	randomNumber, _ := cryptorand.Int(cryptorand.Reader, big.NewInt(number))
	return int(randomNumber.Int64())
}

func RemoveDuplicateStringFromSlice(arr []string) []string {
	// Create a map to track unique elements
	unique := make(map[string]bool)
	var result []string

	for _, num := range arr {
		if _, found := unique[num]; !found {
			unique[num] = true
			result = append(result, num)
		}
	}

	return result
}

func Contains[T comparable](slice []T, element T) bool {
	for _, v := range slice {

		fmt.Println(v, element)

		if v == element {
			return true
		}
	}
	return false
}

func isNumberString(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}
