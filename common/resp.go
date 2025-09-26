package common

type Response struct {
	Data any    `json:"data,inline,omitempty"`
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
