package validation

import (
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
