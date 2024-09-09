package types

import "errors"

type ExitError struct {
	code int
	err  error
}

func NewExitError(code int, err error) *ExitError {
	if err == nil {
		err = errors.New("exit error")
	}
	return &ExitError{
		code: code,
		err:  err,
	}
}

func (e *ExitError) Error() string {
	return e.err.Error()
}

func (e *ExitError) ExitCode() int {
	return e.code
}
