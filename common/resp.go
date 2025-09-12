package common

type Response struct {
	Data any    `json:"data,inline,omitempty"`
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}
