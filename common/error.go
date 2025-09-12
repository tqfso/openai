package common

import "fmt"

type Error struct {
	Code int
	Msg  string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s[%X]", e.Msg, e.Code)
}

func GetErrorCode(err error) int {

	if err == nil {
		return Success
	}

	result, ok := err.(*Error)
	if !ok {
		return Success
	}

	return result.Code
}

func IsErrorCode(err error, code int) bool {
	result, ok := err.(*Error)
	if !ok {
		return false
	}

	return result.Code == code
}
