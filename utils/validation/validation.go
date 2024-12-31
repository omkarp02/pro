package validation

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/utils/errutil"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	validate := validator.New(validator.WithRequiredStructEnabled())
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
