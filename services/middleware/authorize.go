package middleware

import (
	"fmt"

	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/utils"
	"github.com/omkarp02/pro/utils/constant"
	"github.com/omkarp02/pro/utils/errutil"
)

func IsAdmin() router.Handler {

	return func(c router.Context) error {

		role := c.GetDecodedData().Role

		if utils.Contains(role, constant.ROLE_ADMIN) {
			return c.Next()
		}

		return errutil.UnAuthorized("You don't have enough access")
	}
}

func IsOwner() router.Handler {

	return func(c router.Context) error {

		role := c.GetDecodedData().Role

		fmt.Println(role, "<<<<< hellothisisrole<<<<<<")

		if utils.Contains(role, constant.ROLE_OWNER) {
			return c.Next()
		}

		return errutil.UnAuthorized("You don't have enough access")
	}
}
