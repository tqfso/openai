package proxy

import "net/http"

type Error struct {
	Code    int    `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

type ResponseError struct {
	Data Error `json:"error"`
}

func (r ResponseError) Error() string {
	return r.Data.Message
}

func NewResponseError(code int, message string) *ResponseError {
	return &ResponseError{
		Data: Error{
			Code:    code,
			Type:    http.StatusText(code),
			Message: message,
		},
	}
}
