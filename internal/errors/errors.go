package errors

import (
	stdErrors "errors"
	"fmt"
)

type ErrorCode string

const (
	ErrDBLoginAlredyExists    ErrorCode = "login already exists"
	ErrHDRUnmarshallingJSON   ErrorCode = "error unmarshalling json"
	ErrUCWrongLoginOrPassword ErrorCode = "wrong login/password"
	ErrNotAuthenticated       ErrorCode = "not authenticated"
	ErrNoDataFound            ErrorCode = "no data found"
	ErrNotUniqueToken         ErrorCode = "session token already exists"
	ErrDB                     ErrorCode = "some error in storage layer"
	ErrSessionExpired         ErrorCode = "session token is expired"
)

type domainError struct {
	// We define our domainError struct, which is composed of error
	error
	errorCode ErrorCode
}

func (e domainError) Error() string {
	return fmt.Sprintf("%s: %s", e.error.Error(), e.errorCode)
}

func Unwrap(err error) error {
	var dErr domainError
	if stdErrors.As(err, &dErr) {
		return stdErrors.Unwrap(dErr.error)
	}

	return stdErrors.Unwrap(err)
}

func Code(err error) ErrorCode {
	if err == nil {
		return ""
	}

	var dErr domainError
	if stdErrors.As(err, &dErr) {
		return dErr.errorCode
	}

	return ""
}

func NewDomainError(errorCode ErrorCode, format string, args ...interface{}) error {
	return domainError{
		error:     fmt.Errorf(format, args...),
		errorCode: errorCode,
	}
}

func WrapIntoDomainError(err error, errorCode ErrorCode, msg string) error {
	return domainError{
		error:     fmt.Errorf("%s: [%w]", msg, err),
		errorCode: errorCode,
	}
}
