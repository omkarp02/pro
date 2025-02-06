package validation

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/utils/errutil"
)

var validStatuses = map[string]bool{
	"active":   true,
	"inActive": true,
}

func validateStatus(fl validator.FieldLevel) bool {
	// Check if the value exists in the validStatuses map
	status := fl.Field().String()
	return validStatuses[status]
}

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("validStatus", validateStatus)

	return &Validator{
		validate: validate,
	}
}

func (v *Validator) ValidateBody(c router.Context, objType interface{}) error {
	obj := objType

	if err := c.Bind(obj); err != nil {
		return errutil.InvalidReqData()
	}

	if err := v.validate.Struct(obj); err != nil {
		return errutil.HandleValidationError(err)
	}

	return nil
}

func (v *Validator) ValidateParam(c router.Context, key string) (string, error) {
	param := c.Params(key)
	if len(param) == 0 {
		return "", errutil.InvalidReqData()
	}

	return param, nil
}

func (v *Validator) ValidateParams(c router.Context, objType interface{}) error {

	t := reflect.TypeOf(objType)
	if t.Kind() != reflect.Ptr {
		return errutil.InternalServerError("failed to parse interface must be a pointer")
	}

	if err := c.QueryParser(objType); err != nil {
		return errutil.InternalServerError("failed to parse query")
	}

	if err := v.validate.Struct(objType); err != nil {
		return errutil.HandleValidationError(err)
	}

	return nil
}
