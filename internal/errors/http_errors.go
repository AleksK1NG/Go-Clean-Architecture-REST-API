package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var (
	BadRequest              = errors.New("Bad request")
	NoSuchUser              = errors.New("Пользователь с данным именем не существует")
	WrongCredentials        = errors.New("Wrong Credentials")
	NotFound                = errors.New("Не найдено")
	SessionStoreError       = errors.New("Store session error")
	NoSuchSession           = errors.New("Session cookie not found")
	DeleteSessionError      = errors.New("Error while deleting session")
	Unauthorized            = errors.New("Unauthorized")
	SessionTypeAssertionErr = errors.New("Error assert to session type")
	UserTypeAssertionErr    = errors.New("Error assert to user type")
	PermissionsError        = errors.New("У вас недостаточно прав")
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
)

type RestErr interface {
	Message() string
	Status() int
	Error() string
	Causes() []interface{}
}

type restErr struct {
	ErrMessage string        `json:"message,omitempty"`
	ErrStatus  int           `json:"status,omitempty"`
	ErrError   string        `json:"error,omitempty"`
	ErrCauses  []interface{} `json:"causes,omitempty"`
}

func (e restErr) Error() string {
	return fmt.Sprintf("message: %s - status: %d - errors: %s - causes: %v",
		e.ErrMessage, e.ErrStatus, e.ErrError, e.ErrCauses)
}

func (e restErr) Message() string {
	return e.ErrMessage
}

func (e restErr) Status() int {
	return e.ErrStatus
}

func (e restErr) Causes() []interface{} {
	return e.ErrCauses
}

func NewRestError(message string, status int, err string, causes []interface{}) RestErr {
	return restErr{
		ErrMessage: message,
		ErrStatus:  status,
		ErrError:   err,
		ErrCauses:  causes,
	}
}

func NewRestErrorFromBytes(bytes []byte) (RestErr, error) {
	var apiErr restErr
	if err := json.Unmarshal(bytes, &apiErr); err != nil {
		return nil, errors.New("invalid json")
	}
	return apiErr, nil
}

func NewBadRequestError(message string) RestErr {
	return restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusBadRequest,
		ErrError:   "Bad Request",
	}
}

func NewNotFoundError(message string) RestErr {
	return restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusNotFound,
		ErrError:   "Not found",
	}
}

func NewUnauthorizedError(message string) RestErr {
	return restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusUnauthorized,
		ErrError:   "Unauthorized",
	}
}

func NewInternalServerError(message string, err error) RestErr {
	result := restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusInternalServerError,
		ErrError:   "Internal server error",
	}
	if err != nil {
		result.ErrCauses = append(result.ErrCauses, err.Error())
	}
	return result
}

func ParseErrors(err error) RestErr {
	if strings.Contains(err.Error(), "SQLSTATE") {
		return NewBadRequestError("")
	}

	if strings.Contains(err.Error(), "Field validation") {
		return NewBadRequestError(err.Error())
	}

	if err != nil {
		switch err {
		case BadRequest:
			return NewBadRequestError("")
		case NotFound:
			return NewNotFoundError(err.Error())
		case NoSuchUser:
			return NewBadRequestError(err.Error())
		case WrongCredentials:
			return NewBadRequestError(err.Error())
		case Unauthorized:
			return NewUnauthorizedError(err.Error())
		case NotRequiredFields:
			return NewBadRequestError(err.Error())
		case BadQueryParams:
			return NewBadRequestError(err.Error())
		case SessionStoreError:
			return NewBadRequestError("")
		case NoSuchSession:
			return NewUnauthorizedError(Unauthorized.Error())
		case DeleteSessionError:
			return NewUnauthorizedError(Unauthorized.Error())
		case SessionTypeAssertionErr:
			return NewUnauthorizedError(Unauthorized.Error())
		case PermissionsError:
			return NewUnauthorizedError(err.Error())
		case UserTypeAssertionErr:
			return NewInternalServerError("", nil)
		case ExpiredCSRFError:
			return NewUnauthorizedError(Unauthorized.Error())
		case WrongCSRFtoken:
			return NewUnauthorizedError(Unauthorized.Error())
		case CSRFNotPresented:
			return NewUnauthorizedError(Unauthorized.Error())
		case HashingError:
			return NewInternalServerError("", nil)
		case ErrorRequestValidation:
			return NewBadRequestError(BadRequest.Error())
		case ApiAnswerEmptyResult:
			return NewInternalServerError("", nil)
		case ApiResponseStatusNotOK:
			return NewInternalServerError("", nil)
		default:
			return NewInternalServerError("", nil)
		}
	}
	return NewInternalServerError("", nil)
}

// Error response
func ErrorResponse(err error) (int, interface{}) {
	return ParseErrors(err).Status(), ParseErrors(err)
}
