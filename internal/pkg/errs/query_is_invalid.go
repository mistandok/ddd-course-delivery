package errs

import (
	"errors"
	"fmt"
)

var ErrQueryIsInvalid = errors.New("query is invalid")

type QueryIsInvalidError struct {
	QueryName string
	Cause     error
}

func NewQueryIsInvalidErrorWithCause(queryName string, cause error) *QueryIsInvalidError {
	return &QueryIsInvalidError{
		QueryName: queryName,
		Cause:     cause,
	}
}

func NewQueryIsInvalidError(queryName string) *QueryIsInvalidError {
	return &QueryIsInvalidError{
		QueryName: queryName,
	}
}

func (e *QueryIsInvalidError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: query is: %s (cause: %v)",
			ErrQueryIsInvalid, e.QueryName, e.Cause)
	}
	return fmt.Sprintf("%s: %s", ErrQueryIsInvalid, e.QueryName)
}

func (e *QueryIsInvalidError) Unwrap() error {
	return ErrQueryIsInvalid
}
