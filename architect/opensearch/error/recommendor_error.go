package error

import (
	"net/http"
)

const (
	BadRequestErrorType ErrType = iota + 1
	UnauthorizedErrorType
	NotFoundErrorType
	InternalServerErrorType
)

const (
	BadRequestErrorCode         = "ERR_RECOM_BAD_REQUEST_ERROR"
	UnauthorizedActionErrorCode = "ERR_RECOM_UNAUTHORIZED_ACTION_ERROR"
	NotFoundErrorCode           = "ERR_RECOM_NOT_FOUND_ERROR"
	InternalServerErrorCode     = "ERR_RECOM_INTERNAL_SERVER_ERROR"
)

type ErrType uint

func (errorType ErrType) New(errorCode ErrCode, errorMessage string, statusCode int) *RECOMError {
	return &RECOMError{
		ErrType: errorType,
		ErrResponse: ErrResponse{
			ErrCode:    errorCode,
			ErrMessage: errorMessage,
		},
		StatusCode: statusCode,
	}
}

type ErrCode string

type ErrResponse struct {
	ErrCode    ErrCode `json:"error_code"`
	ErrMessage string  `json:"error_message"`
}

type RECOMError struct {
	StatusCode  int
	ErrType     ErrType
	ErrResponse ErrResponse
}

var ErrRecordNotFound = NotFoundErrorType.New(NotFoundErrorCode, "record not found", http.StatusNotFound)
var ErrURLQueryParamNotFound = BadRequestErrorType.New(BadRequestErrorCode, "url query parameter not found", http.StatusBadRequest)
var ErrUnAuthorizedAction = UnauthorizedErrorType.New(UnauthorizedActionErrorCode, "not authorized for action ", http.StatusUnauthorized)

func BadRequestErrorFunc(errorMessage string) *RECOMError {
	return BadRequestErrorType.New(BadRequestErrorCode, errorMessage, http.StatusBadRequest)
}

func InternalServerErrorFunc(errorMessage string) *RECOMError {
	return InternalServerErrorType.New(InternalServerErrorCode, errorMessage, http.StatusInternalServerError)
}
