package errors

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/go.net/context"
	"net/http"
	"strings"
)

var (
	BadRequest              = errors.New("Bad request")
	NoSuchUser              = errors.New("User not found")
	WrongCredentials        = errors.New("Wrong Credentials")
	NotFound                = errors.New("Not Found")
	SessionStoreError       = errors.New("Store session error")
	NoSuchSession           = errors.New("Session cookie not found")
	DeleteSessionError      = errors.New("Error while deleting session")
	Unauthorized            = errors.New("Unauthorized")
	Forbidden               = errors.New("Forbidden")
	SessionTypeAssertionErr = errors.New("Error assert to session type")
	UserTypeAssertionErr    = errors.New("Error assert to user type")
	PermissionsError        = errors.New("У вас недостаточно прав")
	PermissionDenied        = errors.New("Permission Denied")
	ExpiredCSRFError        = errors.New("Expired CSRF token")
	WrongCSRFtoken          = errors.New("Wrong CSRF token")
	CSRFNotPresented        = errors.New("CSRF not presented")
	HashingError            = errors.New("Ошибка хеширования")
	NotRequiredFields       = errors.New("No such required fields")
	ErrorRequestValidation  = errors.New("Ошибка валидации")
	BadQueryParams          = errors.New("Invalid query params")
	ApiResponseStatusNotOK  = errors.New("Ошибка овтета внешнего API")
	ApiAnswerEmptyResult    = errors.New("Некорректный ответ внешнего API")
	InternalServerError     = errors.New("Internal Server Error")
	RequestTimeoutError     = errors.New("Request Timeout")
	ExistsEmailError        = errors.New("User with given email already exists")
	InvalidJWTToken         = errors.New("Invalid JWT token")
	InvalidJWTClaims        = errors.New("Invalid JWT claims")
	InvalidJWTSignMethod    = errors.New("Invalid JWT sign method")
)

type RestErr interface {
	Status() int
	Error() string
	Causes() interface{}
}

type restErr struct {
	ErrStatus int         `json:"status,omitempty"`
	ErrError  string      `json:"error,omitempty"`
	ErrCauses interface{} `json:"-"`
}

func (e restErr) Error() string {
	return fmt.Sprintf("status: %d - errors: %s - causes: %v", e.ErrStatus, e.ErrError, e.ErrCauses)
}

func (e restErr) Status() int {
	return e.ErrStatus
}

func (e restErr) Causes() interface{} {
	return e.ErrCauses
}

func NewRestError(status int, err string, causes interface{}) RestErr {
	return restErr{
		ErrStatus: status,
		ErrError:  err,
		ErrCauses: causes,
	}
}

func NewRestErrorFromBytes(bytes []byte) (RestErr, error) {
	var apiErr restErr
	if err := json.Unmarshal(bytes, &apiErr); err != nil {
		return nil, errors.New("invalid json")
	}
	return apiErr, nil
}

func NewBadRequestError(causes interface{}) RestErr {
	return restErr{
		ErrStatus: http.StatusBadRequest,
		ErrError:  BadRequest.Error(),
		ErrCauses: causes,
	}
}

func NewNotFoundError(causes interface{}) RestErr {
	return restErr{
		ErrStatus: http.StatusNotFound,
		ErrError:  NotFound.Error(),
		ErrCauses: causes,
	}
}

func NewUnauthorizedError(causes interface{}) RestErr {
	return restErr{
		ErrStatus: http.StatusUnauthorized,
		ErrError:  Unauthorized.Error(),
		ErrCauses: causes,
	}
}

func NewForbiddenError(causes interface{}) RestErr {
	return restErr{
		ErrStatus: http.StatusForbidden,
		ErrError:  Forbidden.Error(),
		ErrCauses: causes,
	}
}

func NewInternalServerError(causes interface{}) RestErr {
	result := restErr{
		ErrStatus: http.StatusInternalServerError,
		ErrError:  InternalServerError.Error(),
		ErrCauses: causes,
	}

	return result
}

func ParseErrors(err error) RestErr {

	if strings.Contains(err.Error(), "SQLSTATE") {
		return parseSqlErrors(err)
	}

	if strings.Contains(err.Error(), "Field validation") {
		return parseValidatorError(err)
	}

	if err == sql.ErrNoRows {
		return NewRestError(http.StatusNotFound, NotFound.Error(), err)
	}

	if err == context.DeadlineExceeded {
		return NewRestError(http.StatusRequestTimeout, RequestTimeoutError.Error(), err)
	}

	if strings.Contains(err.Error(), "Unmarshal") {
		return NewRestError(http.StatusBadRequest, BadRequest.Error(), err)
	}

	if strings.Contains(err.Error(), "UUID") {
		return NewRestError(http.StatusBadRequest, err.Error(), err)
	}

	if strings.Contains(strings.ToLower(err.Error()), "cookie") {
		return NewRestError(http.StatusUnauthorized, Unauthorized.Error(), err)
	}

	if strings.Contains(strings.ToLower(err.Error()), "token") {
		return NewRestError(http.StatusUnauthorized, Unauthorized.Error(), err)
	}

	if strings.Contains(strings.ToLower(err.Error()), "bcrypt") {
		return NewRestError(http.StatusBadRequest, BadRequest.Error(), err)
	}

	restErr, ok := err.(RestErr)
	if ok {
		return ParseRestErrors(restErr)
	}

	return NewInternalServerError(err)
}

func parseSqlErrors(err error) RestErr {
	if strings.Contains(err.Error(), "23505") {
		return NewRestError(http.StatusBadRequest, ExistsEmailError.Error(), err)
	}

	return NewRestError(http.StatusBadRequest, BadRequest.Error(), err)
}

func parseValidatorError(err error) RestErr {
	if strings.Contains(err.Error(), "Password") {
		return NewRestError(http.StatusBadRequest, "Invalid password, min length 6", err)
	}

	if strings.Contains(err.Error(), "Email") {
		return NewRestError(http.StatusBadRequest, "Invalid email", err)
	}

	return NewRestError(http.StatusBadRequest, BadRequest.Error(), err)
}

func ParseRestErrors(err RestErr) RestErr {
	if err != nil {
		switch err {
		case PermissionDenied:
			return NewUnauthorizedError(err)
		case BadRequest:
			return NewBadRequestError(err)
		case NotFound:
			return NewNotFoundError(err)
		case NoSuchUser:
			return NewBadRequestError(err)
		case WrongCredentials:
			return NewBadRequestError(err)
		case Unauthorized:
			return NewUnauthorizedError(err)
		case NotRequiredFields:
			return NewBadRequestError(err)
		case BadQueryParams:
			return NewBadRequestError(err)
		case SessionStoreError:
			return NewBadRequestError(err)
		case NoSuchSession:
			return NewUnauthorizedError(err)
		case DeleteSessionError:
			return NewUnauthorizedError(err)
		case SessionTypeAssertionErr:
			return NewUnauthorizedError(err)
		case PermissionsError:
			return NewUnauthorizedError(err)
		case UserTypeAssertionErr:
			return NewInternalServerError(err)
		case ExpiredCSRFError:
			return NewUnauthorizedError(err)
		case WrongCSRFtoken:
			return NewUnauthorizedError(err)
		case CSRFNotPresented:
			return NewUnauthorizedError(err)
		case HashingError:
			return NewInternalServerError(err)
		case ErrorRequestValidation:
			return NewBadRequestError(err)
		case ApiAnswerEmptyResult:
			return NewInternalServerError(err)
		case ApiResponseStatusNotOK:
			return NewInternalServerError(err)
		default:
			return NewInternalServerError(err)
		}
	}
	return NewInternalServerError(err)
}

// Error response
func ErrorResponse(err error) (int, interface{}) {
	return ParseErrors(err).Status(), ParseErrors(err)
}
