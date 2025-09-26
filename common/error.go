package common

type Error struct {
	Code int
	Msg  string
}

func (e Error) Error() string {
	return e.Msg
}

func GetErrorCode(err error, def int) int {
	result, ok := err.(*Error)
	if !ok {
		return def
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
