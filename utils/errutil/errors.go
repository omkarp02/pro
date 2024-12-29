package errutil

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type APIError struct {
	StatusCode int    `json:"statusCode"`
	Code       int    `json:"code"`
	Msg        string `json:"message"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("api error => msg : %d => code : %d", e.StatusCode, e.Code)
}

func NewAPIError(statusCode int, code int, msg string) APIError {
	return APIError{
		StatusCode: statusCode,
		Code:       code,
		Msg:        msg,
	}
}

var ErrDocumentNotFound = NewAPIError(http.StatusInternalServerError, 2, "Document Not found")
var ErrDocumentAlreadyExist = NewAPIError(http.StatusInternalServerError, 3, "Document Already Exist")
var ErrDatabase = NewAPIError(http.StatusInternalServerError, 8, "Database Error")

func GenerateError(statusCode int, code int, err error) error {
	return APIError{
		StatusCode: statusCode,
		Code:       code,
		Msg:        err.Error(),
	}
}

func HandleValidationError(errs error) error {

	var errMsgs []string

	for _, err := range errs.(validator.ValidationErrors) {
		field := err.Field()

		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("Field %s is required field", field))
		case "email":
			errMsgs = append(errMsgs, fmt.Sprintf("Field %s should be a valid email", field))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("Field %s is invalid", field))
		}
	}
	return NewAPIError(http.StatusUnprocessableEntity, 7, strings.Join(errMsgs, ","))
}

func InvalidCredentails() APIError {
	return NewAPIError(http.StatusUnauthorized, 4, "Invalid Credentials")
}

func UnAuthorized(msg string) APIError {
	return NewAPIError(http.StatusUnauthorized, 5, msg)
}

func InternalServerError(msg ...string) APIError {
	m := "Internal Server Error"
	if len(msg) > 0 {
		m = msg[0]
	}
	return NewAPIError(http.StatusInternalServerError, 6, m)
}

func InvalidReqData() APIError {
	return NewAPIError(http.StatusBadRequest, 6, "Invalid Request Data")
}

func AlreadyExist(name string) APIError {
	return NewAPIError(http.StatusInternalServerError, 2, name+"already exist")
}

func StatusBadRequest(msg string) error {
	return NewAPIError(http.StatusBadRequest, 1, msg)
}
