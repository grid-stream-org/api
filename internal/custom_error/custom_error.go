package custom_error

import (
	"fmt"

	"github.com/pkg/errors"
)

type CustomError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"err"`
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s, Error: %s", e.Code, e.Message, errors.WithStack(e.Err))
}

func New(code int, message string, err error) *CustomError {
    if err == nil {
        err = errors.New(message)
    }
	return &CustomError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
